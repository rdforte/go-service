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

// deleteUser removes a user from the system.
func (h userHandler) deleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return validate.NewRequestError(validate.ErrForbidden, http.StatusForbidden)
	}

	userID := claims.Subject

	if err := h.user.Delete(ctx, userID); err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidID):
			return validate.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, user.ErrNotFound):
			return validate.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", userID, err)
		}
	}

	return web.RespondOk(ctx, w)
}
