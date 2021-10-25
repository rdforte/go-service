package web

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	Status string `json:"status"`
}

type Handler func(res http.ResponseWriter, req *http.Request)

type route struct {
	method          map[string]Handler
	notFoundHandler Handler
}

func (r *route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if handler, ok := r.method[req.Method]; !ok {
		r.notFoundHandler(res, req)
	} else {
		handler(res, req)
	}
}

type App struct {
	*http.ServeMux
	notFoundHandler Handler
	router          map[string]*route // 1st key = url, 2nd key = http method
}

func createDefaultNotFoundHandler() Handler {
	return func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(Status{"NOT FOUND"})
	}
}

func NewApp() *App {
	return &App{
		http.NewServeMux(),
		createDefaultNotFoundHandler(),
		make(map[string]*route),
	}
}

func (a *App) SetupRoute(path, method string, handler Handler) {
	if _, ok := a.router[path]; !ok {
		a.router[path] = &route{
			method:          make(map[string]Handler),
			notFoundHandler: a.notFoundHandler,
		}
		a.router[path].method[method] = handler
		a.Handle(path, a.router[path])
	}

	a.router[path].method[method] = handler
}

func (a *App) Get(path string, handler Handler) {
	a.SetupRoute(path, http.MethodGet, handler)
}

func (a *App) Post(path string, handler Handler) {
	a.SetupRoute(path, http.MethodPost, handler)
}

func (a *App) Patch(path string, handler Handler) {
	a.SetupRoute(path, http.MethodPatch, handler)
}

func (a *App) Delete(path string, handler Handler) {
	a.SetupRoute(path, http.MethodDelete, handler)
}

// Call before all routes you would like to set custom error messages
func (a *App) NotFound(handler func(res http.ResponseWriter, req *http.Request)) {
	a.notFoundHandler = handler
}
