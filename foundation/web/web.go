// Package web contains a small web famework extension.
package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	/**
	  Currently we are embedding gorilla/mux as the main mux and hiding this away through the App struct.
	  This gives us the ability to easily switch over the mux by just changing the embedding
	  and some minor api logic in the foundation/web layer.
	*/
	"github.com/google/uuid"
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
	mw       []Middleware
}

// NewApp creates an App vaue that handles a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		mux.NewRouter(),
		shutdown,
		mw,
	}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// handleReq is the main method we use to build are App based networking handlers.
func (a *App) handleReq(path string, group string, httpMethod string, handler Handler, mw ...Middleware) {
	p := path
	if group != "" {
		p = fmt.Sprintf("/%s%s", group, path)
	}

	/**
	First wrap handler specific middleware.
	Second add the app specific middleware after, Which results in the app specific middleware being executed
	prior to handler specific middleware.
	*/
	handler = wrapMiddleWare(mw, handler)
	handler = wrapMiddleWare(a.mw, handler)

	a.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {

		// PRE CODE PROCESSING
		ctx := r.Context()

		v := Values{
			TracedID: uuid.New().String(),
			Now:      time.Now(),
		}

		ctx = context.WithValue(ctx, key, &v)

		// Call the wrapped handler.
		if err := handler(ctx, w, r); err != nil {
			// The Error should not reach the outer most handler. If it does it means we have an issue with our system.
			// In this case signal shutdown.
			a.SignalShutdown()
		}

		// POST CODE PROCESSING
	}).Methods(httpMethod)
}

/**
Get handler for handling all http GET requests.
Calls to this handler are used for Reading data.
*/
func (a *App) Get(path string, group string, handler Handler, mw ...Middleware) {
	a.handleReq(path, group, http.MethodGet, handler, mw...)
}

/**
Post handler for handling all http POST requests.
Calls to this handler are used for writing data.
*/
func (a *App) Post(path string, group string, handler Handler, mw ...Middleware) {
	a.handleReq(path, group, http.MethodPost, handler, mw...)
}

/**
Patch handler for handling all http PATCH requests.
Calls to this handler are used to update/modify an existing resource.
*/
func (a *App) Patch(path string, group string, handler Handler, mw ...Middleware) {
	a.handleReq(path, group, http.MethodPatch, handler, mw...)
}

/**
Put handler for handling all http PUT requests.
Calls to this handler are used for replacing/overriding an existing resource.
*/
func (a *App) Put(path string, group string, handler Handler, mw ...Middleware) {
	a.handleReq(path, group, http.MethodPut, handler, mw...)
}

/**
Delete handler for handling all http DELETE requests.
Calls to this handler are used for deleting a resource.
*/
func (a *App) Delete(path string, group string, handler Handler, mw ...Middleware) {
	a.handleReq(path, group, http.MethodDelete, handler, mw...)
}
