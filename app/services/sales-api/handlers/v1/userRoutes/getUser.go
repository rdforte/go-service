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

// QueryByID returns a user by its ID.
func (h userHandler) getUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return validate.NewRequestError(errors.New("attempted action is not allowed"), http.StatusForbidden)
	}

	userID := claims.Subject

	usr, err := h.user.QueryByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidID):
			return validate.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, user.ErrNotFound):
			return validate.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", userID, err)
		}
	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}
