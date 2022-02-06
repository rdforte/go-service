package main

import (
	"bytes"
	"context"
	"encoding/json"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rdforte/go-service/app/services/sales-api/handlers"
	"github.com/rdforte/go-service/business/sys/auth"
	"github.com/rdforte/go-service/foundation/keystore"
	"github.com/spf13/viper"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/**
* Keep TODO's at the top of the file ie:
* 1. Figure out timeouts for http service.
 */

var build = "develop"

func main() {
	log, err := initLogger("SALES-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// =========================================================================================================
	// GOMAXPROCS

	// Set the current number of threads for the service
	// based on what is available by either the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}

	// =========================================================================================================
	// CONFIGURATION

	type Config struct {
		Version struct {
			SVN  string `yaml:"svn"`
			Desc string `yaml:"desc"`
		}
		Web struct {
			ReadTimeout     int    `yaml:"readTimeout"`
			WriteTimeout    int    `yaml:"writeTimeout"`
			IdleTimeout     int    `yaml:"idleTimeout"`
			ShutdownTimeout int    `yaml:"shutdownTimeout"`
			APIHost         string `yaml:"apiHost"`
			DebugHost       string `yaml:"debugHost"`
		}
		Auth struct {
			KeysFolder string `yaml:"keysFolder"`
			ActiveKID  string `yaml:"activeKID"`
		}
	}

	cfg := &Config{}

	viper.AddConfigPath("../../config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("error reading config", "ERROR", err)
	}

	viper.Set("version.svn", build)

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Errorw("err unmarshaling confing", "ERROR", err)
	}

	// format the config to be pretty when logged so we can easily identify the configuration of the service.
	src, err := json.Marshal(cfg)
	if err != nil {
		log.Errorw("err marshaling config", "ERROR", err)
	}
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, src, "", " "); err != nil {
		log.Errorw("can not format config", "ERROR", err)
	}

	log.Infow("Setting up service with config", "config", dst.String())

	// =========================================================================================================
	// Authentication
	ks, err := keystore.NewFS()
	if err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	auth, err := auth.New(cfg.Auth.ActiveKID, ks)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// =========================================================================================================
	// APP STARTING

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	// set the build number when identifying metrics in expvar
	expvar.NewString("build").Set(build)

	// =========================================================================================================
	// APP STARTING
	log.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	/** The Debug function returns a mux to listen and serve on for all the debug
	related endpoints. This includes the standard library endpoints.
	*/
	debugMux := handlers.DebugMux(build, log)

	// start the service listening for debug requests.
	// not concerned about shutting this down with load shedding.
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// =========================================================================================================

	// Macke a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Construct the mux for the API calls.
	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
		Auth:     auth,
	})

	// Construct a server to service the requests against a mux
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  time.Duration(cfg.Web.ReadTimeout),
		WriteTimeout: time.Duration(cfg.Web.WriteTimeout),
		IdleTimeout:  time.Duration(cfg.Web.IdleTimeout),
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// Make a channel to listen for erros coming from the listener.
	// Use a buffered channel so the goroutine can exit if we don't collect this error.
	// When we shutdown the server ListenAndServe can return straight away because we have a buffer channel of 1.
	/** If it was unbuffered then then sender and reciever of the channel need to be in sync and because we are
	running the shutdown case in the select with a timeout then that wont be recieving and therefore will block the return
	of the ListenAndServe when its time to shutdown. Because the channel is buffered and there is some buffer space available
	then the ListenAndServe can send the error to the serverErrors channel without there needing to be a receiver waiting to
	receive on the other side therefore allowing the shutdown of listenAndServe to start as soon as we signal the shutdown of the
	server.
	*/
	serverErrors := make(chan error, 1)

	// Start the service listening for api requests.
	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Web.ShutdownTimeout))
		defer cancel()

		// Asking listener to shutdown and shed load.
		// Shutdown is blocking and will take the context of the timeout.
		if err := api.Shutdown(ctx); err != nil {
			api.Close() // if the shutdown times out then close server manually.
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

// Construct the application logger.
func initLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}
