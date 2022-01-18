package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	// log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =========================================================================
	// CONFIGURATION

	type Config struct {
		Version struct {
			SVN  string `json:"svn"`
			Desc string `json:"desc"`
		}
		// web struct {
		// 	readTimeout     int
		// 	writeTimeout    int
		// 	idleTimeout     int
		// 	shutdownTimeout int
		// 	apiHost         string
		// 	debugHost       string
		// }
	}

	conf := &Config{}

	viper.AddConfigPath("./app/config")
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		log.Infow("error reading config", "error", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Infow("err unmarshaling confing", "error", err)
	}
	// ver := viper.Get("version")
	fmt.Println("--->", *&conf.Version.SVN)
	// const prefix = "SALES"
	// help, err := conf.ParseOSArgs(prefix, &cfg)
	// if err != nil {
	// 	if errors.Is(err, conf.ErrHelpWanted) {
	// 		fmt.Println(help)
	// 		return nil
	// 	}
	// 	return fmt.Errorf("parsing config: %w", err)
	// }

	// =========================================================================
	// APP STARTING

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	// out, err := conf.String(&cfg)
	// if err != nil {
	// 	return fmt.Errorf("generating config for output: %w", err)
	// }
	// log.Infow("startup", "config", out)

	// =========================================================================
	// APP STARTING
	log.Infow("startup", "status", "debug router started", "host")

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
