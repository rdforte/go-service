package userRoutes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rdforte/go-service/business/core/user"
	"github.com/rdforte/go-service/business/sys/auth"
	"github.com/rdforte/go-service/foundation/web"
)

// The fields we expect the client to send when they signup
type decodeUser struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// Create adds a new user to the system.
func (h userHandler) signUp(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var du decodeUser

	if err := web.Decode(r, &du); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	// construct user for saving in database.
	nu := user.NewUser{
		Name:            du.Name,
		Email:           du.Email,
		Password:        du.Password,
		PasswordConfirm: du.PasswordConfirm,
		Roles:           []string{auth.RoleUser},
	}

	usr, err := h.user.Create(ctx, nu, v.Now)
	if err != nil {
		return fmt.Errorf("user[%+v]: %w", &usr, err)
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID,
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	tok, err := h.auth.GenerateToken(claims)
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  cookieKey,
		Value: tok,
	})

	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	statusCode := http.StatusOK

	return web.Respond(ctx, w, status, statusCode)
}
