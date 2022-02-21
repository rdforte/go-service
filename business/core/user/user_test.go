package user_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

	tl.Describe("Working with User records")
	{
		tl.It("should be able to handle a single user")

		ctx := context.Background()
		now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)

		nu := user.NewUser{
			Name:            "Master Yoda",
			Email:           "masteryoda@1lbladfaulbahnaldjbjbadanl23456.com",
			Roles:           []string{auth.RoleAdmin},
			Password:        "gophers",
			PasswordConfirm: "gophers",
		}

		// Create User.
		usr, err := core.Create(ctx, nu, now)
		if err != nil {
			tl.Failed("Should be able to create user", err)
		}
		tl.Success("Should be able to create user")

		// Query User by ID.
		saved, err := core.QueryByID(ctx, usr.ID)
		if err != nil {
			tl.Failed("Should be able to retrieve user by ID", err)
		}
		tl.Success("Should be able to retrive user by ID")

		// Compare user created to user queried.
		if diff := cmp.Diff(usr, saved); diff != "" {
			tl.Failed("Should get back the same user", fmt.Errorf("created user does not match queried user"))
		}
		tl.Success("Should get back the same user")

		// Update user.
		upd := user.UpdateUser{
			Name:  dbtest.StringPointer("Luke Skywalker"),
			Email: dbtest.StringPointer("lukeskywalker@anfuaofuanofuaofafnaofuafaofa.com.au"),
		}

		if err := core.Update(ctx, usr.ID, upd, now); err != nil {
			tl.Failed("Should be able to update user", err)
		}
		tl.Success("Should be able to update user")

		// Query user by Email.
		saved, err = core.QueryByEmail(ctx, *upd.Email)
		if err != nil {
			tl.Failed("Should be able to retrieve user by email", err)
		}
		tl.Success("Should be able to retrieve user by email")

		// Check if updated Email is correct.
		if saved.Email != *upd.Email {
			tl.Failed("Should have been able to update email",
				fmt.Errorf("Emails do not match: [%s : %s]", saved.Email, *upd.Email))
		}
		tl.Success("Should have been able to update email")

		// Check if updated Name is correct.
		if saved.Name != *upd.Name {
			tl.Failed("Should have been able to update name",
				fmt.Errorf("Names do not match: [%s : %s]", saved.Name, *upd.Name))
		}
		tl.Success("Should have been able to update name")

		// Delete user.
		if err := core.Delete(ctx, usr.ID); err != nil {
			tl.Failed("Should be able to delete user", err)
		}
		tl.Success("Should be able to delete user")

		// Retrieve Deleted user expecting user to not be in db.
		if _, err := core.QueryByID(ctx, usr.ID); !errors.Is(err, user.ErrNotFound) {
			tl.Failed("Should not be able to retrieve user", err)
		}
		tl.Success("Should not be ablt to retrieve user")
	}
}
