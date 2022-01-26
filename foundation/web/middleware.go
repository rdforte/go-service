package web

/**
Middleware is a function designed to run some code before and/or after
another handler. It is designed to remove boilerplate or other concerns not
direct to any given Handler. The Middlware promises to wrap and call the handler provided
while returning a new handler.
*/
type Middleware func(Handler) Handler

/**
wrapMiddleWare creates a new Handler by wrapping middleware around a final
Handler. The middleware's Handlers will be executed by requests in the order
they are provided
*/
func wrapMiddleWare(mw []Middleware, handler Handler) Handler {

	/**
	Loop backwards through the middleware invoking each middleware.
	Replace the handler with the new wrapped handler. Looping backwards
	ensures that the first middleware of the slice is the first to be executed
	by the request.
	*/
	for i := len(mw) - 1; i >= 0; i-- {
		mWare := mw[i]
		if mWare != nil {
			handler = mWare(handler)
		}
	}
	return handler
}
