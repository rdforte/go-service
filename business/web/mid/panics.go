package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rdforte/go-service/foundation/web"
)

/**
Panics recovers from panics and converts the panic to an error so it is reported in metrics
and handled in Errors.
*/
func Panics() web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached to the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// Defer a function to recover from a panic and set the err return variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {

					// get the stack trace for the logs to help with debugging the panic
					trace := debug.Stack()

					// Stack trace will be provided.
					// Because we are in a defer statement we have to assign the error to the err variable
					// as a means of returning the error to the calling function.
					err = fmt.Errorf("PANIC [%v] TRACE [%s]", rec, string(trace))
				}
			}()

			// Proceed with calling the next handler. If the next handler panics we will catch it in the
			// above defer statement.
			return handler(ctx, w, r)
		}
		return h
	}
	return m
}
