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

	"github.com/rdforte/go-service/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/rdforte/go-service/app/services/sales-api/handlers/v1/testgrp"
	"github.com/rdforte/go-service/business/sys/auth"
	"github.com/rdforte/go-service/business/web/mid"
	"github.com/rdforte/go-service/foundation/web"
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
	Auth     *auth.Auth
}

// APIMux constructs an http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {

	// set up the web app with app specific middleware
	// Panics must always be at the end so that it is the first middleware to be called around the
	// handler in case there is a panic within the handler we can handle this.
	r := web.NewApp(cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	// Load the routes for the different versions of the API.
	v1(r, cfg)

	return r
}

func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}

	app.Get("/test", version, tgh.Test)
	app.Get("/testauth", version, tgh.Test, mid.Authenticate(cfg.Auth), mid.Authorize("ADMIN"))
}
