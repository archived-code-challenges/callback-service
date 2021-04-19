package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opencensus.io/trace"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Handler is the signature used by all application handlers in this service.
type Handler func(context.Context, http.ResponseWriter, *http.Request)

// Middleware is a function designed to run some code before and/or after
// another Handler. It is designed to remove boilerplate or other concerns not
// direct to any given Handler.
type Middleware func(Handler) Handler

// App is the entrypoint into our application and what controls the context of
// each request. Feel free to add any configuration data/logic on this type.
type App struct {
	log *log.Logger
	mux *chi.Mux
	mw  []Middleware
}

// NewApp constructs an App to handle a set of routes. Any Middleware provided
// will be ran for every request.
func NewApp(log *log.Logger, mw ...Middleware) *App {
	return &App{
		log: log,
		mux: chi.NewRouter(),
		mw:  mw,
	}
}

// wrapMiddleware creates a new handler by wrapping middleware around a final
// handler. The middlewares' Handlers will be executed by requests in the order
// they are provided.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}

// Values carries information about each request.
type Values struct {
	TraceID    string
	StatusCode int
	Start      time.Time
}

// Handle associates a handler function with an HTTP Method and URL pattern.
//
// It converts our custom handler type to the std lib Handler type. It captures
// errors from the handler and serves them to the client in a uniform way.
func (a *App) Handle(method, url string, h Handler) {

	// wrap the application's middleware around this endpoint's handler.
	h = wrapMiddleware(a.mw, h)

	// Create a function that conforms to the std lib definition of a handler.
	// This is the first thing that will be executed when this route is called.
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "internal.platform.web")
		defer span.End()

		// Create a Values struct to record state for the request. Store the
		// address in the request's context so it is sent down the call chain.
		v := Values{
			TraceID: span.SpanContext().TraceID.String(),
			Start:   time.Now(),
		}
		ctx = context.WithValue(ctx, KeyValues, &v)

		h(ctx, w, r)
	}

	a.mux.MethodFunc(method, url, fn)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
