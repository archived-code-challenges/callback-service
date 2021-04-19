package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.opencensus.io/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	serverObjectsEndpointURL = "/objects/"
)

var (
	CallbackSelfDeleteTime = 30 * time.Second
)

type Callback struct {
	ID        int64 `gorm:"primary_key;type:bigserial" json:"id"`
	Online    bool  `gorm:"not null" json:"online"`
	Timestamp int64 `gorm:"type:bigint;not null" json:"timestamp"`
}

type callbackService struct {
	CallbackService
}

// CallbackService defines a set of methods to be used when dealing when a callback is received.
type CallbackService interface {
	// Upsert will start a goroutine function to fetch the status of any Callback provided and then
	// call the database.Upsert method.
	// Any error that happens here, will be logged, won't be returned.
	Upsert(context.Context, []Callback) error

	// Status fetches the callback status from the callback-client service and fills the fetched Online status
	Status(context.Context, int64) (Callback, error)

	CallbackDB
}

// CallbackDB defines how the service interacts with the database.
type CallbackDB interface {
	// Upsert a slice of Callbacks on the database. Will ignore any ID conflict.
	// This function is configured to delete any callback object that is upserted 30 seconds after
	// its upsertion if its timestamp is from 30 seconds ago.
	Upsert(context.Context, []Callback) error
}

func NewCallbackService(db *gorm.DB, callbackServiceURL string, log *log.Logger) CallbackService {
	return &callbackService{
		CallbackService: &callbackValidator{
			CallbackDB: &callbackGorm{db},
			serviceURL: callbackServiceURL,
			log:        log,
		},
	}
}

type callbackValidator struct {
	CallbackDB

	serviceURL string
	log        *log.Logger
	ctx        context.Context
}

// setTimestamp sets timestamp to now. It does not return any errors.
func (cv *callbackValidator) setTimestamp(c *Callback) {
	c.Timestamp = time.Now().Unix()
}

func (cv *callbackValidator) Upsert(ctx context.Context, cs []Callback) error {
	ctx, span := trace.StartSpan(ctx, "models.callbackValidator.Upsert")
	defer span.End()

	errChan := make(chan error)
	done := make(chan bool)
	var wg sync.WaitGroup

	for _, c := range cs {
		wg.Add(1)
		go func(c Callback) { // The c argument is used to capture the loop variable at the moment is used.

			// Use client to fetch callback status
			callback, err := cv.Status(ctx, c.ID)
			if err != nil {
				errChan <- err
			}

			if callback.Online {
				// Run callback validators
				cv.setTimestamp(&callback)

				// Upsert callback on the database
				err = cv.CallbackDB.Upsert(ctx, []Callback{callback})
				if err != nil {
					errChan <- err
				}
			}

			wg.Done()
		}(c)
	}

	// goroutine to wait until WaitGroup is done
	go func() {
		wg.Wait()
		close(done)
	}()

	// goroutine to wait until either WaitGroup is done or an error is received through the channel
	go func() {
		select {
		case <-done:
			break
		case err := <-errChan:
			cv.log.Printf("upsert_error: %v", err)
		}
	}()

	return nil
}

func (cv *callbackValidator) Status(ctx context.Context, id int64) (Callback, error) {
	_, span := trace.StartSpan(ctx, "models.callbackValidator.Status")
	defer span.End()

	buildURL := fmt.Sprintf("%s%s%s", cv.serviceURL, serverObjectsEndpointURL, strconv.FormatInt(id, 10))
	req, err := http.NewRequest(http.MethodGet, buildURL, nil)
	if err != nil {
		return Callback{}, fmt.Errorf("models: building request %w", err)
	}

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		return Callback{}, fmt.Errorf("models: sending http request %w", err)
	}
	defer resp.Body.Close()

	var nc Callback
	if err := json.NewDecoder(resp.Body).Decode(&nc); err != nil {
		return Callback{}, ErrInvalidJSONInput
	}

	return nc, nil
}

type callbackGorm struct {
	db *gorm.DB
}

func (cg *callbackGorm) Upsert(ctx context.Context, cs []Callback) error {
	ctx, span := trace.StartSpan(ctx, "callback.Database.Upsert")
	defer span.End()
	cg.db.WithContext(ctx)

	err := cg.db.Debug().Clauses(clause.OnConflict{
		UpdateAll: true, // Update everything on ID conflict.
	}).Create(&cs).Error

	if err != nil {
		return fmt.Errorf("models: couldn't update callback %w", err)
	}

	// Slice containing the IDs of the Callback objects just created.
	var bulkDeleteIDs []int64
	for _, c := range cs {
		bulkDeleteIDs = append(bulkDeleteIDs, c.ID)
	}

	// Function that will run exactly X seconds after the Callback objects were created and will target them.
	// Deleting them if their timestamp has not been updated.
	time.AfterFunc(CallbackSelfDeleteTime, func() {
		cg.db.Debug().Where("timestamp <= ?", time.Now().Add(-CallbackSelfDeleteTime).Unix()).Delete(Callback{}, bulkDeleteIDs)
	})

	return nil
}
