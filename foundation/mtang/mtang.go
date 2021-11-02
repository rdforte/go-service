package mtang

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type Status struct {
	Status string `json:"status"`
}

type Handler func(ctx Context, res http.ResponseWriter, req *http.Request)

type Route struct {
	method    map[string]Handler
	regPath   *regexp.Regexp
	pathParam *pathParam
}

type Router struct {
	*http.ServeMux
	notFoundHandler Handler
	router          map[string]*Route // 1st key = url, 2nd key = http method
}

func createDefaultNotFoundHandler() Handler {
	return func(ctx Context, res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(Status{"NOT FOUND"})
	}
}

type apiHandler struct {
	app *Router
}

func (a *apiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.String())
	for _, route := range a.app.router {
		if match := route.regPath.Match([]byte(req.URL.String())); match {
			if handler, ok := route.method[req.Method]; ok {
				ctx := createNewReqCtx(req)
				if len(route.pathParam.keys) > 0 {
					pathChunks := route.pathParam.pathSegmentRgx.FindAllString(req.URL.String(), -1)
					fmt.Println(pathChunks)
					for i, p := range route.pathParam.positions {
						param := pathChunks[p]
						key := route.pathParam.keys[i]
						ctx.params[key] = param
					}

				}
				handler(ctx, res, req)
				return
			}
		}
	}
	a.app.notFoundHandler(Context{}, res, req)
}

func NewRouter() *Router {
	app := &Router{
		http.NewServeMux(),
		createDefaultNotFoundHandler(),
		make(map[string]*Route),
	}
	app.Handle("/", &apiHandler{app})
	return app
}

func (r *Router) SetupRoute(path, method string, handler Handler) {
	if _, ok := r.router[path]; !ok {
		r.router[path] = buildRoute(path)
		r.router[path].method[method] = handler
	}

	r.router[path].method[method] = handler
}

func (a *Router) Get(path string, handler Handler) {
	a.SetupRoute(path, http.MethodGet, handler)
}

func (a *Router) Post(path string, handler Handler) {
	a.SetupRoute(path, http.MethodPost, handler)
}

func (r *Router) Patch(path string, handler Handler) {
	r.SetupRoute(path, http.MethodPatch, handler)
}

func (r *Router) Delete(path string, handler Handler) {
	r.SetupRoute(path, http.MethodDelete, handler)
}

// Call before all routes you would like to set custom error messages
func (r *Router) NotFound(handler func(ctx Context, res http.ResponseWriter, req *http.Request)) {
	r.notFoundHandler = handler
}

type pathChunk struct {
	position []int
	pathType string //path | param
}

type pathParam struct {
	pathSegmentRgx regexp.Regexp
	positions      []int    // the index in the path for which the param is located
	keys           []string // the keys associated with the param
}
