package middleware

import (
	"context"
	"expvar"
	"net/http"
	"runtime"

	"go.opencensus.io/trace"

	"github.com/noelruault/go-callback-service/internal/web"
)

// m contains the global program counters for the application.
var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

// Metrics updates program counters.
func Metrics() web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		// Wrap this handler around the next one provided.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			ctx, span := trace.StartSpan(ctx, "internal.middleware.Metrics")
			defer span.End()

			before(ctx, w, r)

			// Increment the request counter.
			m.req.Add(1)

			// Update the count for the number of active goroutines every 100 requests.
			if m.req.Value()%100 == 0 {
				m.gr.Set(int64(runtime.NumGoroutine()))
			}
		}

		return h
	}

	return f
}
