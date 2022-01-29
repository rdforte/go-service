package metrics

import (
	"context"
	"expvar"
)

/**
This holds the single instance of the metrics value needed for collecting metrics. The expvar package
is already based on the singleton for the different metrics that are registered with the package so
there isn't much choice here.
*/
var m *metrics

// ===========================================================================================================

/**
Metrics represents the set of metrics we gather. These fields are safe to be accessed concurrently thanks
to expvar. No extra abstraction is required.
*/
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

/**
init construcs the metrics value that will be used to capture metrics. The metrics value is stored in a
package level variable since everything inside of expvar is registered as a singleton. The use of init will
make sure this initialization only happens once.
*/
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

// ===========================================================================================================

// ctxKeyMetric represents the type of value for the context key.
type ctxKey int

// key is how metric values are stored/retrieved.
const key ctxKey = 1

// Set sets the metrics data into the context.
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

// AddGoroutines increments the goroutines metric by 1.
func AddGoroutines(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.goroutines.Value()%100 == 0 {
			v.goroutines.Add(1)
		}
	}
}

// AddRequests increments the request metric by 1.
func AddRequests(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)
	}
}

// AddErrors increments the errors metric by 1.
func AddErrors(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
	}
}

// AddPanics increments the panics metric by 1.
func AddPanics(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
	}
}
