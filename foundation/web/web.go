package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type Status struct {
	Status string `json:"status"`
}

type Handler func(res http.ResponseWriter, req *http.Request)

type Route struct {
	method          map[string]Handler
	notFoundHandler Handler
}

func (r *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if handler, ok := r.method[req.Method]; !ok {
		r.notFoundHandler(res, req)
	} else {
		handler(res, req)
	}
}

type App struct {
	*http.ServeMux
	notFoundHandler Handler
	Router          map[string]*Route // 1st key = url, 2nd key = http method
}

func createDefaultNotFoundHandler() Handler {
	return func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(Status{"NOT FOUND"})
	}
}

type apiHandler struct {
	app *App
}

func (a *apiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for key, route := range a.app.Router {
		regPath := regexp.MustCompile(key)
		if match := regPath.Match([]byte(req.URL.String())); match {
			if handler, ok := route.method[req.Method]; ok {
				handler(res, req)
				return
			}
		}
		fmt.Println(key)
	}
	a.app.notFoundHandler(res, req)
}

func NewApp() *App {
	app := &App{
		http.NewServeMux(),
		createDefaultNotFoundHandler(),
		make(map[string]*Route),
	}
	app.Handle("/", &apiHandler{app})
	return app
}

func (a *App) SetupRoute(path, method string, handler Handler) {
	if _, ok := a.Router[path]; !ok {
		a.Router[path] = &Route{
			method:          make(map[string]Handler),
			notFoundHandler: a.notFoundHandler,
		}
		a.Router[path].method[method] = handler
		// a.Handle(path, a.Router[path])
	}

	a.Router[path].method[method] = handler
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
