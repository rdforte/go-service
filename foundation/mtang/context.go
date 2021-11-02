package mtang

import (
	"context"
	"net/http"
	"net/url"
)

type Context struct {
	context.Context
	params  map[string]string
	queries url.Values
}

func (ctx *Context) GetQuery(key string) string {
	return ctx.queries.Get(key)
}

func (ctx *Context) GetParam(key string) string {
	return ctx.params[key]
}

func createNewReqCtx(req *http.Request) Context {
	return Context{
		req.Context(),
		map[string]string{},
		req.URL.Query(),
	}
}
