package userRoutes

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rdforte/go-service/business/core/user"
	"github.com/rdforte/go-service/business/sys/auth"
	"github.com/rdforte/go-service/business/sys/validate"
	"github.com/rdforte/go-service/foundation/web"
)

// update updates a user in the system.
func (h userHandler) updateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return validate.NewRequestError(validate.ErrForbidden, http.StatusForbidden)
	}

	var upd user.UpdateUser
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	userID := claims.Subject

	if err := h.user.Update(ctx, userID, upd, v.Now); err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidID):
			return validate.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, user.ErrNotFound):
			return validate.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] User[%+v]: %w", userID, &upd, err)
		}
	}

	return web.RespondOk(ctx, w)
}
