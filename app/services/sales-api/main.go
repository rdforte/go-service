package main

import (
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "expvar"

	"github.com/rdforte/go-service/app/services/sales-api/handlers"
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
	// =========================================================================
	// GOMAXPROCS

	// Set the current number of threads for the service
	// based on what is available by either the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}

	// =========================================================================
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
			ApiHost         string `yaml:"apiHost"`
			DebugHost       string `yaml:"debugHost"`
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

	// =========================================================================
	// APP STARTING

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	// set the build number when identifying metrics in expvar
	expvar.NewString("build").Set(build)

	// =========================================================================
	// APP STARTING
	log.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	/** The Debug function returns a mux to listen and serve on for all the debug
	related endpoints. This includes the standard library endpoints.
	*/
	debugMux := handlers.DebugStandardLibraryMux()

	// start the service listening for debug requests.
	// not concerned about shutting this down with load shedding.
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// =========================================================================

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

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
