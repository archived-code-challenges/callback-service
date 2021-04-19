package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.opencensus.io/trace"
	"gorm.io/gorm"

	"github.com/noelruault/go-callback-service/internal/models"
	"github.com/noelruault/go-callback-service/internal/web"
)

// Callback defines all of the handlers related to products. It holds the application state needed by the handler methods.
type Callback struct {
	csvc models.CallbackService

	log *log.Logger
}

// NewCallbacks creates a new Callback controller.
func NewCallbacks(csvc models.CallbackService, log *log.Logger) *Callback {

	return &Callback{
		csvc: csvc,
		log:  log,
	}
}

type callbackRequest struct {
	Objects []int64 `json:"object_ids"`
}

// Retrieve finds a single product identified by an ID in the request URL.
func (c *Callback) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(ctx, "handlers.Callback.Handle")
	defer span.End()

	var cr callbackRequest
	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		web.RespondError(ctx, w, ErrInvalidJSONInput, http.StatusBadRequest)
		return
	}

	// Remove any possible duplicate value on the request IDs
	ids := removeDuplicateValues(cr.Objects)

	// Build a slice of Callbacks that will be upserted
	var callbackList []models.Callback
	for _, id := range ids {
		callbackList = append(callbackList, models.Callback{ID: id})
	}

	// Send Callback slice to Upsert method
	err := c.csvc.Upsert(ctx, callbackList)
	if err != nil {
		web.RespondError(ctx, w, err, http.StatusNotAcceptable)
		return
	}

	web.Respond(ctx, w, struct{}{}, http.StatusOK)
}

// Check provides support for orchestration health checks.
type Check struct {
	db  *gorm.DB
	log *log.Logger
}

// Health validates the service is healthy and ready to accept requests.
func (c *Check) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	_, span := trace.StartSpan(ctx, "handlers.Check.Health")
	defer span.End()

	var health struct {
		Status string `json:"status"`
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if pinger, ok := c.db.ConnPool.(interface{ Ping() error }); ok {
		err := pinger.Ping()
		if err != nil {
			health.Status = "database couldn't be reached"
			w.WriteHeader(http.StatusInternalServerError)
		}
		if !ok {
			health.Status = "database driver/type not supported"
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	health.Status = "ok"

	data, err := json.Marshal(health)
	if err != nil {
		c.log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		c.log.Println("error writing result", err)
	}

	return
}
