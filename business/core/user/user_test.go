package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/rdforte/go-service/business/core/user"
	"github.com/rdforte/go-service/business/data/dbtest"
	"github.com/rdforte/go-service/business/sys/auth"
	"github.com/rdforte/go-service/foundation/logger"
)

var dbc = dbtest.DBContainer{
	Image: "postgres:14-alpine",
	Port:  "5432",
	Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
}

func TestUser(t *testing.T) {
	testID := 0
	tl := logger.NewTestLog(t, testID)

	log, db, teardown := dbtest.NewUnit(t, dbc)
	t.Cleanup(teardown)

	core := user.NewCore(log, db)

	tl.Describe("Working with User records.")
	{
		tl.It("should be able to handle a single user")

		ctx := context.Background()
		now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)

		nu := user.NewUser{
			Name:            "Ryan Forte",
			Email:           "ryan@testing123456.com",
			Roles:           []string{auth.RoleAdmin},
			Password:        "gophers",
			PasswordConfirm: "gophers",
		}

		usr, err := core.Create(ctx, nu, now)
		if err != nil {
			tl.Failed("Should be able to create user", err)
		}
		tl.Success("Should be able to create user")
		t.Log(usr)
	}
}
