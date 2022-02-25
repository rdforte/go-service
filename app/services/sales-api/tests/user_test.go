package user_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rdforte/go-service/app/services/sales-api/handlers"
	"github.com/rdforte/go-service/business/core/user"
	"github.com/rdforte/go-service/business/data/dbtest"
	"github.com/rdforte/go-service/business/sys/auth"
	"github.com/rdforte/go-service/foundation/logger"
	"github.com/rdforte/go-service/foundation/web"
)

// TODO - include tests for deleteUser handler
// TODO - include tests for signUp handler
// TODO - include tests for updateUser handler
// TODO - include tests for unhappy path login handler
// TODO - include tests for unhappy path getUser handler

// UserTests holds methods for each user subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered.
type UserTests struct {
	app       http.Handler
	userToken string
	tl        *logger.TestLogger
}

// TestUsers is the entry point for testing user management functions.
func TestUsers(t *testing.T) {
	test := dbtest.NewIntegration(
		t,
		dbtest.DBContainer{
			Image: "postgres:14-alpine",
			Port:  "5432",
			Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
		},
	)
	t.Cleanup(test.Teardown)

	tl := logger.NewTestLog(t)
	tl.Describe("User Handlers")

	shutdown := make(chan os.Signal, 1)
	tests := UserTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.Auth,
			DB:       test.DB,
		}),
		userToken: test.Token("user@example.com", "gophers"),
		tl:        tl,
	}

	t.Run("Login200", tests.loginSuccess)
	t.Run("Login200", tests.getUserSuccess)

}

// loginSuccess tests the happy path for a successful user login.
func (ut *UserTests) loginSuccess(t *testing.T) {
	ut.tl.It("Should be able to login a user successfully")

	r := httptest.NewRequest(http.MethodPost, "/v1/user/login", nil)
	w := httptest.NewRecorder()

	r.SetBasicAuth("user@example.com", "gophers")

	ut.app.ServeHTTP(w, r)

	// Sets correct status code.
	if w.Code != http.StatusOK {
		ut.tl.Failed("should return status 200", fmt.Errorf("Status [%d]", w.Code))
	}
	ut.tl.Success("should return status 200")

	// Returns the correct json payload.
	b := &web.OK{}
	if err := json.Unmarshal(w.Body.Bytes(), b); err != nil {
		ut.tl.Failed("should be able to see response body with status : OK", err)
	}
	if b.Status != "OK" {
		ut.tl.Failed("should be able to see response body with status : OK", fmt.Errorf("did not return status OK"))
	}
	ut.tl.Success("should be able to see response body with status : OK")

	// Sets the token in the cookies.
	if w.Result().Cookies()[0].Value != ut.userToken {
		ut.tl.Failed("should be able to set token in cookies", fmt.Errorf("cookies: %v", w.Result().Cookies()))
	}
	ut.tl.Success("should be able to set token in cookies")
}

// getUserSuccess tests the happy path for getting a users details.
func (ut *UserTests) getUserSuccess(t *testing.T) {
	ut.tl.It("Should be able to retrieve a users details if authenticated")

	r := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
	w := httptest.NewRecorder()

	r.AddCookie(&http.Cookie{
		Name:  "xra789klate",
		Value: ut.userToken,
	})

	ut.app.ServeHTTP(w, r)

	// Sets correct status code.
	if w.Code != http.StatusOK {
		ut.tl.Failed("should return status 200", fmt.Errorf("Status [%d]", w.Code))
	}
	ut.tl.Success("should return status 200")

	// Returns the correct user details.
	usr := &user.User{}
	if err := json.Unmarshal(w.Body.Bytes(), usr); err != nil {
		ut.tl.Failed("should be able to see response body with status : OK", err)
	}

	usrDate, _ := time.Parse("2006-01-02 00:00:00 +0000 UTC", "2019-03-24 00:00:00 +0000 UTC")

	if usr.ID != "45b5fbd3-755f-4379-8f07-a58d4a30fa2f" ||
		len(usr.PasswordHash) != 0 || // password should not be included in response.
		usr.Name != "User Gopher" ||
		usr.Email != "user@example.com" ||
		len(usr.Roles) != 1 ||
		usr.Roles[0] != auth.RoleUser ||
		usr.DateCreated != usrDate ||
		usr.DateUpdated != usrDate {
		ut.tl.Failed("should return the correct user details", fmt.Errorf("user: %+v", usr))
	}
	ut.tl.Success("should return the correct user details")
}
