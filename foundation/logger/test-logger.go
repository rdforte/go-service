package logger

import "testing"

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// TestLog constructs
type testLogger struct {
	t  *testing.T
	id int
}

// NewTestLog constructs a new Test Logger for handling Success and Failure logs in tests.
func NewTestLog(t *testing.T, testID int) *testLogger {
	return &testLogger{
		t,
		testID,
	}
}

// Success is a wrapper function for Logf that helps with formating the Log.
func (tl *testLogger) Success(message string) {
	tl.t.Logf("\t%s\tTest %d:\t%s.", Success, tl.id, message)

}

// Failed is a wrapper function for Fatalf that helps with formating the Log.
func (tl *testLogger) Failed(message string, err error) {
	tl.t.Fatalf("\t%s\tTest %d:\t%s : %s.", Failed, tl.id, message, err)
}
