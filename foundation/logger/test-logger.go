package logger

import "testing"

// Success and failure markers.
const (
	success = "\u2705"
	failed  = "\u274c"
	light   = "🚨"
)

// TestLog constructs
type TestLogger struct {
	t *testing.T
}

// NewTestLog constructs a new Test Logger for handling Success and Failure logs in tests.
func NewTestLog(t *testing.T) *TestLogger {
	return &TestLogger{
		t,
	}
}

// Success is a wrapper function for Logf that helps with formating the Log.
func (tl *TestLogger) Success(message string) {
	tl.t.Logf("\033[32m\t%s\tTest:\t%s.\033[0m", success, message)

}

// Failed is a wrapper function for Fatalf that helps with formating the Log.
func (tl *TestLogger) Failed(message string, err error) {
	tl.t.Fatalf("\033[31m\t%s\tTest:\t%s : %s.\033[0m", failed, message, err)
}

// Describe describes what it is you are testing. Use this as your first log to group your assertions.
func (tl *TestLogger) Describe(message string) {
	tl.t.Logf("%s%s \033[35m%s\033[0m %s%s\n\n", light, light, message, light, light)
}

// It is used for adding clarity to your assertions ie: `It should authenticate a single user`.
// Use this after the Describe.
// All messages are prefixed with It.
func (tl *TestLogger) It(message string) {
	tl.t.Logf("\033[33m\tTest:\tIt %s\n\n\033[0m", message)
}
