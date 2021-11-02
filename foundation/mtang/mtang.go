package mtang

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	Status string `json:"status"`
}

type Handler func(ctx Context, res http.ResponseWriter, req *http.Request)

type Router struct {
	*http.ServeMux
	notFoundHandler Handler
	routes          map[string]*route // 1st key = route, 2nd key = route information associated with the url
}

func createDefaultNotFoundHandler() Handler {
	return func(ctx Context, res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(Status{"NOT FOUND"})
	}
}

type entryPoint struct {
	router *Router
}

func (e *entryPoint) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, route := range e.router.routes {
		if match := route.regPath.Match([]byte(req.URL.String())); match {
			if handler, ok := route.method[req.Method]; ok {
				ctx := createNewReqCtx(req)
				if len(route.pathParam.keys) > 0 {
					pathChunks := route.pathParam.pathSegmentRgx.FindAllString(req.URL.String(), -1)
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
	e.router.notFoundHandler(Context{}, res, req)
}

// Creates a new Router
func NewRouter() *Router {
	router := &Router{
		http.NewServeMux(),
		createDefaultNotFoundHandler(),
		make(map[string]*route),
	}
	router.Handle("/", &entryPoint{router})
	return router
}

func (r *Router) SetupRoute(path, method string, handler Handler) {
	if _, ok := r.routes[path]; !ok {
		r.routes[path] = buildRoute(path)
		r.routes[path].method[method] = handler
	}

	r.routes[path].method[method] = handler
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

// overrides the default NotFound handler
func (r *Router) NotFound(handler func(ctx Context, res http.ResponseWriter, req *http.Request)) {
	r.notFoundHandler = handler
}
