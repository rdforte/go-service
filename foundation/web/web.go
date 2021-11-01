package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
)

type Status struct {
	Status string `json:"status"`
}

type Handler func(ctx Context, res http.ResponseWriter, req *http.Request)

type Route struct {
	method          map[string]Handler
	notFoundHandler Handler
	regPath         *regexp.Regexp
	pathParam       *pathParam
}

func (r *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if handler, ok := r.method[req.Method]; !ok {
		r.notFoundHandler(Context{}, res, req)
	} else {
		handler(Context{}, res, req)
	}
}

type App struct {
	*http.ServeMux
	notFoundHandler Handler
	Router          map[string]*Route // 1st key = url, 2nd key = http method
}

func createDefaultNotFoundHandler() Handler {
	return func(ctx Context, res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(Status{"NOT FOUND"})
	}
}

type Context struct {
	context.Context
	Params map[string]string
}

func createNewReqCtx(ctx context.Context) Context {
	return Context{
		ctx,
		map[string]string{},
	}
}

type apiHandler struct {
	app *App
}

func (a *apiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, route := range a.app.Router {
		if match := route.regPath.Match([]byte(req.URL.String())); match {
			if handler, ok := route.method[req.Method]; ok {
				ctx := createNewReqCtx(req.Context())
				if len(route.pathParam.keys) > 0 {
					pathChunks := route.pathParam.pathSegmentRgx.FindAllString(req.URL.String(), -1)
					fmt.Println(pathChunks)
					for i, p := range route.pathParam.positions {
						param := pathChunks[p]
						key := route.pathParam.keys[i]
						ctx.Params[key] = param
					}

				}
				handler(ctx, res, req)
				return
			}
		}
	}
	a.app.notFoundHandler(Context{}, res, req)
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
	rgx, pp := regBuilder(path)
	fmt.Println(pp)

	if _, ok := a.Router[path]; !ok {
		regPath := rgx
		a.Router[path] = &Route{
			method:          make(map[string]Handler),
			notFoundHandler: a.notFoundHandler,
			regPath:         regPath,
			pathParam:       pp,
		}
		a.Router[path].method[method] = handler
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
func (a *App) NotFound(handler func(ctx Context, res http.ResponseWriter, req *http.Request)) {
	a.notFoundHandler = handler
}

func regBuilder(path string) (rgx *regexp.Regexp, pp *pathParam) {
	fmt.Println(len(path))
	paramRegex := regexp.MustCompile(`\/:[\d\w]+([^/])`)
	pathRegex := regexp.MustCompile(`\/[^:][\d\w]+([^/])`)
	params := paramRegex.FindAllIndex([]byte(path), -1)
	paths := pathRegex.FindAllIndex([]byte(path), -1)
	fmt.Println(paths)
	fmt.Println(params)

	pathChunks := []pathChunk{}
	for _, val := range paths {
		chunk := pathChunk{
			val,
			"path",
		}
		pathChunks = append(pathChunks, chunk)
		// fmt.Println(absolutePath)
	}
	for _, val := range params {
		chunk := pathChunk{
			val,
			"param",
		}
		pathChunks = append(pathChunks, chunk)
		// fmt.Println(absolutePath)
	}
	sort.Slice(pathChunks, func(i, j int) bool {
		return pathChunks[i].position[0] < pathChunks[j].position[0]
	})

	regPath := "^"
	pp = &pathParam{
		*regexp.MustCompile(`[^\/:][\w\d-_]+`),
		[]int{},
		[]string{},
	}
	for i, val := range pathChunks {
		if val.pathType == "path" {
			regPath += path[val.position[0]:val.position[1]]
		} else if val.pathType == "param" {
			regPath += `/[\d\w_-]+`
			pp.positions = append(pp.positions, i)
			key := regexp.MustCompile(`[^\/:][\w\d-_]+`).Find([]byte(path[val.position[0]:val.position[1]]))
			pp.keys = append(pp.keys, string(key))
		}
		if i == len(pathChunks)-1 {
			regPath += "$"
		}
	}
	fmt.Println(regPath)
	return regexp.MustCompile(regPath), pp
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
