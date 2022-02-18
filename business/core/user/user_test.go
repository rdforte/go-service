package user_test

import (
	"testing"

	"github.com/rdforte/go-service/business/data/dbtest"
)

var dbc = dbtest.DBContainer{
	Image: "postgres:14-alpine",
	Port:  "5432",
	Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
}

func TestUser(t *testing.T) {
	// testID := 0
	// tl := logger.NewTestLog(t, testID)

	// log, db, teardown := dbtest.NewUnit(t, dbc)
	// t.Cleanup(teardown)

	// core := user.NewCore(log, db)

	t.Log("Given the need to work with User records.")
	{

	}
}
