/**
Package handlers contains the full set of handler functions and routes
supported by the web api.
*/
package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/gorilla/mux"
	"github.com/rdforte/go-service/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/rdforte/go-service/app/services/sales-api/handlers/v1/testgrp"
	"go.uber.org/zap"
)

/**
debugStandardLibraryMux registers all the debug routes from the standard library
into a new mux bypassing the use of the DefaultServeMux. Using the the DefaultServeMux
would be a security risk since a dependency could inject a handler into our service
without us knowing about it.
*/
func debugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

/**
DebugMux registers all the debug standard library routes and then custom debug
application routes for the service. This bypasses the use of the DefaultServerMux.
Using the DefaultServerMux would be a security risk since a dependency could inject
a handler into our service without us knowing it.
*/
func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := debugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}

	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// APIMuxConfig contains all the mandatory systems required by the handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux constructs an http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) http.Handler {
	r := mux.NewRouter()

	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}

	r.HandleFunc("/v1/test", tgh.Test).Methods("GET")

	return r
}
