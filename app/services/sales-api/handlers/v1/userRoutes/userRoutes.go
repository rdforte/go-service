// Package userRoutes maintains the group of handlers for user access.
package userRoutes

import (
	"github.com/rdforte/go-service/business/core/user"
	"github.com/rdforte/go-service/business/sys/auth"
	"github.com/rdforte/go-service/business/web/mid"
	"github.com/rdforte/go-service/foundation/web"
)

type userHandler struct {
	user user.Core
	auth *auth.Auth
}

// CreateUserV1Routes is a function responsible for setting up all the V1 User routes.
func CreateUserV1Routes(app *web.App, user user.Core, auth *auth.Auth) {
	// Create User Handler
	usrHandler := userHandler{
		user,
		auth,
	}

	authenticate := mid.Authenticate(auth)

	// User Routes
	app.Post("/user/login", "v1", usrHandler.login)
	app.Post("/user/signup", "v1", usrHandler.signUp)

	// User Routes (Authenticated)
	app.Get("/user", "v1", usrHandler.getUser, authenticate)
	app.Patch("/user", "v1", usrHandler.updateUser, authenticate)
	app.Delete("/user", "v1", usrHandler.deleteUser, authenticate)
}
