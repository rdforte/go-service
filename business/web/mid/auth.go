package mid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rdforte/go-service/business/sys/auth"
	"github.com/rdforte/go-service/business/sys/validate"
	"github.com/rdforte/go-service/foundation/web"
)

// cookieKey is the key for the token when set in the cookies.
var cookieKey = "xra789klate"

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(a *auth.Auth) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Expecting: bearer <token>
			// authStr := r.Header.Get("authorization")
			c, err := r.Cookie(cookieKey)
			if err != nil {
				return validate.NewRequestError(err, http.StatusUnauthorized)
			}

			// Parse the authorization header
			// parts := strings.Split(c.Value, " ")
			// if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			// 	err := errors.New("expecting authorization header format: bearer <token>")
			// 	return validate.NewRequestError(err, http.StatusUnauthorized)
			// }

			// Validate that the token is signed by us
			// claims, err := a.ValidateToken(parts[1])
			claims, err := a.ValidateToken(c.Value)
			if err != nil {
				return validate.NewRequestError(err, http.StatusUnauthorized)
			}

			ctx = auth.SetClaims(ctx, claims)

			return handler(ctx, w, r)
		}
		return h
	}
	return m
}

// Authorize validates that an authenticated user ahs at least one role from a specified list.
// This method constructs the actual function that is used
func Authorize(roles ...string) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			claims, err := auth.GetClaims(ctx)
			if err != nil {
				return validate.NewRequestError(
					fmt.Errorf("you are not authorized for that action, no claims"),
					http.StatusForbidden,
				)
			}

			if !claims.Authorized(roles...) {
				return validate.NewRequestError(
					fmt.Errorf("you are not authorized for the actions, claims[%v] roles[%v]", claims.Roles, roles),
					http.StatusForbidden,
				)
			}

			return handler(ctx, w, r)
		}
		return h
	}
	return m
}
