package mid

import (
	"context"
	"net/http"

	"github.com/rdforte/go-service/business/sys/validate"
	"github.com/rdforte/go-service/foundation/web"
	"go.uber.org/zap"
)

/**
Errors handles errors coming out of the call chain. It detects normal application errors which
are used to respond to the client in a uniform way. Unexpected errors (status >= 500) are logged.
*/
func Errors(log *zap.SugaredLogger) web.Middleware {

	// This the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Creates the handler that will be attached to the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service to be shutdown gracefully.
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}

			// Call the inner handler.
			if err := handler(ctx, w, r); err != nil {

				// Log the error.
				log.Errorw("ERROR", "traceid", v.TracedID, "ERROR", err)

				// Build out the error response.
				var er validate.ErrorResponse
				var status int

				switch act := validate.Cause(err).(type) {
				case validate.FieldErrors:
					er = validate.ErrorResponse{
						Error:  "data validation error",
						Fields: act.Error(),
					}
					status = http.StatusBadRequest

				case *validate.RequestError:
					er = validate.ErrorResponse{
						Error: act.Error(),
					}
					status = act.Status

				default:
					// default is a non trusted error so return status 500.
					er = validate.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				// Respond with the error back to the client.
				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err then we need to return it back to the base handler
				// to shutdown the service.
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

			// The error has been handled so we can stop propogating it.
			return nil
		}
		return h
	}
	return m
}
