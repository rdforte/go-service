// Package web contains a small web famework extension.
package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	/**
	  Currently we are embedding gorilla/mux as the main mux and hiding this away through the App struct.
	  This gives us the ability to easily switch over the mux by just changing the embedding
	  and some minor api logic in the foundation/web layer.
	*/
	"github.com/gorilla/mux"
)

// A Handler is a type that handles an http request within our App.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

/**
App is the entrypoint into our application ans what configures our context
object for each of our http handlers. Feel free to add any configuration
data/logic on this App struct.
*/
type App struct {
	*mux.Router
	shutdown chan os.Signal
}

// NewApp creates an App vaue that handles a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App {
	return &App{
		mux.NewRouter(),
		shutdown,
	}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// handleReq is the main method we use to build are App based networking handlers.
func (a *App) handleReq(path string, group string, httpMethod string, handler Handler) {
	p := path
	if group != "" {
		p = fmt.Sprintf("/%s%s", group, path)
	}

	a.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {

		// PRE CODE PROCESSING
		ctx := r.Context()

		if err := handler(ctx, w, r); err != nil {
			// ERROR HANDLING
		}

		// POST CODE PROCESSING
	}).Methods(httpMethod)
}

/**
Get handler for handling all http GET requests.
Calls to this handler are used for Reading data.
*/
func (a *App) Get(path string, group string, handler Handler) {
	a.handleReq(path, group, http.MethodGet, handler)
}

/**
Post handler for handling all http POST requests.
Calls to this handler are used for writing data.
*/
func (a *App) Post(path string, group string, handler Handler) {
	a.handleReq(path, group, http.MethodPost, handler)
}

/**
Patch handler for handling all http PATCH requests.
Calls to this handler are used to update/modify an existing resource.
*/
func (a *App) Patch(path string, group string, handler Handler) {
	a.handleReq(path, group, http.MethodPatch, handler)
}

/**
Put handler for handling all http PUT requests.
Calls to this handler are used for replacing/overriding an existing resource.
*/
func (a *App) Put(path string, group string, handler Handler) {
	a.handleReq(path, group, http.MethodPut, handler)
}

/**
Delete handler for handling all http DELETE requests.
Calls to this handler are used for deleting a resource.
*/
func (a *App) Delete(path string, group string, handler Handler) {
	a.handleReq(path, group, http.MethodPut, handler)
}
