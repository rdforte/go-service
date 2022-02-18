package logger

import "testing"

// Success and failure markers.
const (
	success = "\u2705"
	failed  = "\u274c"
	light   = "🚨"
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
	tl.t.Logf("\033[32m\t%s\tTest %d:\t%s.\033[0m", success, tl.id, message)

}

// Failed is a wrapper function for Fatalf that helps with formating the Log.
func (tl *testLogger) Failed(message string, err error) {
	tl.t.Fatalf("\033[31m\t%s\tTest %d:\t%s : %s.\033[0m", failed, tl.id, message, err)
}

// Describe describes what it is you are testing. Use this as your first log to group your assertions.
func (tl *testLogger) Describe(message string) {
	tl.t.Log(light, light, " ", message, " ", light, light, "\n")
}

// It is used for adding clarity to your assertions ie: `It should authenticate a single user`.
// Use this after the Describe.
// All messages are prefixed with It.
func (tl *testLogger) It(message string) {
	tl.t.Logf("\033[33m\tTest %d:\tIt %s\n\n\033[0m", tl.id, message)
}
