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
	regPath         *regexp.Regexp
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
	for _, route := range a.app.Router {
		// regPath := regexp.MustCompile(key)
		if match := route.regPath.Match([]byte(req.URL.String())); match {
			if handler, ok := route.method[req.Method]; ok {
				handler(res, req)
				return
			}
		}
		// fmt.Println(key)
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
	// Todo TESTING building routes
	regBuilder(path)

	if _, ok := a.Router[path]; !ok {
		regPath := regexp.MustCompile(path)
		a.Router[path] = &Route{
			method:          make(map[string]Handler),
			notFoundHandler: a.notFoundHandler,
			regPath:         regPath,
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

func regBuilder(path string) *regexp.Regexp {
	fmt.Println(len(path))
	paramRegex := regexp.MustCompile(`\/:[\d\w]+([^/])`)
	pathRegex := regexp.MustCompile(`\/[^:][\d\w]+([^/])`)
	params := paramRegex.FindAllIndex([]byte(path), -1)
	paths := pathRegex.FindAllIndex([]byte(path), -1)
	fmt.Println(paths)
	fmt.Println(params)
	for _, val := range paths {
		absolutePath := path[val[0]:val[1]] // this is the path that we will use to construct the regex
		fmt.Println(path[val[0]:val[1]])
	}
	return paramRegex
}
