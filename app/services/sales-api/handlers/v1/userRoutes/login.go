package userRoutes

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rdforte/go-service/business/core/user"
	"github.com/rdforte/go-service/business/sys/validate"
	"github.com/rdforte/go-service/foundation/web"
)

// cookieKey is the key for the token when set in the cookies.
var cookieKey = "xra789klate"

func (h userHandler) login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return validate.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := h.user.Authenticate(ctx, v.Now, email, pass)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return validate.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, user.ErrAuthenticationFailure):
			return validate.NewRequestError(err, http.StatusUnauthorized)
		default:
			return fmt.Errorf("authenticating: %w", err)
		}
	}

	tok, err := h.auth.GenerateToken(claims)
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  cookieKey,
		Value: tok,
	})

	return web.RespondOk(ctx, w)
}
