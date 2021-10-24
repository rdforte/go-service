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
	method map[string]Handler
}

func (r *route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if handler, ok := r.method[req.Method]; !ok {
		json.NewEncoder(res).Encode(Status{"NOT FOUND"})
	} else {
		handler(res, req)
	}
}

type App struct {
	*http.ServeMux
	notFoundHandler func(res http.ResponseWriter, req *http.Request)
	router          map[string]*route // 1st key = url, 2nd key = http method
}

func NewApp() *App {
	return &App{
		http.NewServeMux(),
		func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusNotFound)
			json.NewEncoder(res).Encode(Status{"NOT FOUND"})
		},
		make(map[string]*route),
	}
}

// func (a *App) HandleRoute(path string, res http.ResponseWriter, req *http.Request) {
// 	if _, ok := a.router[path]; !ok {
// 		a.notFoundHandler(res, req)
// 	} else {
// 		a.Handle(path, &route{})
// 	}
// }

func (a *App) Get(path string, handler Handler) {
	_, ok := a.router[path]

	if !ok {
		a.router[path] = &route{
			method: make(map[string]Handler),
		}
		a.router[path].method[http.MethodGet] = handler
		a.Handle(path, a.router[path])
	}

	a.router[path].method[http.MethodGet] = handler

	// h := func(res http.ResponseWriter, req *http.Request) {
	// 	a.HandleRoute(path, res, req)
	// }
	// a.HandleFunc(path, h)
}

func (a *App) Post(path string, handler Handler) {
	_, ok := a.router[path]

	if !ok {
		a.router[path] = &route{
			method: make(map[string]Handler),
		}
		a.router[path].method[http.MethodGet] = handler
		a.Handle(path, a.router[path])
	}

	a.router[path].method[http.MethodPost] = handler

	// h := func(res http.ResponseWriter, req *http.Request) {
	// 	a.HandleRoute(path, res, req)
	// }
	// a.HandleFunc(path, h)
}

func (a *App) NotFound(handler func(res http.ResponseWriter, req *http.Request)) {
	a.notFoundHandler = handler
}
