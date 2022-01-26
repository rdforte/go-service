package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/rdforte/go-service/foundation/web"
	"go.uber.org/zap"
)

func Logger(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request for the service to be shutdonw gracefuly.
			v, err := web.GetValues(ctx)
			if err != nil {
				return err
			}

			log.Infow("request started", "traceid", v.TracedID, "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr)

			// Call the next handler
			err = handler(ctx, w, r)

			log.Infow("request completed", "traceid", v.TracedID, "method", r.Method, "path", r.URL.Path, "statusCode", v.StatusCode,
				"remoteAddr", r.RemoteAddr, "since", time.Since(v.Now))

			return err
		}
		return h
	}
	return m
}
