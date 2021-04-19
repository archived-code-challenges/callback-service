// +build integration

package handlers_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/noelruault/go-callback-service/internal/handlers"
	"github.com/noelruault/go-callback-service/internal/models"
)

const (
	serverCallbackURL        = "http://localhost:9010"
	serverObjectsEndpointURL = "/objects/"
)

func TestCallback_Integration(t *testing.T) {
	tdb := models.NewTestDatabase(t)
	models.CallbackSelfDeleteTime = 5 * time.Second
	defer func() {
		models.CallbackSelfDeleteTime = 30 * time.Second
		models.CleanupTestDatabase(tdb)
	}()

	testlog := log.New(log.Writer(), "test", 0)
	csvc := models.NewCallbackService(tdb, serverCallbackURL, testlog)
	c := handlers.NewCallbacks(csvc, testlog)

	_, err := http.Get(fmt.Sprintf("%s%s", serverCallbackURL, serverObjectsEndpointURL))
	assert.NoError(t, err, "The service/endpoint is not reachable")

	var cases = []struct {
		name      string
		input     string
		outStatus int
		outJSON   string
	}{
		{
			"notJSON",
			"a dalhd lkald fkjahd lfkjasdlf ",
			http.StatusBadRequest,
			`{"error":"invalid_json","message":"provided input cannot be parsed"}`,
		},
		{
			"ok",
			`{"object_ids": [91,10,78,11,30,40,22,33]}`,
			http.StatusOK,
			`{}`,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewReader([]byte(cs.input)))
			ctx := handlers.NewTestContext()

			c.Handle(ctx, w, r)

			assert.Equal(t, cs.outStatus, w.Result().StatusCode)
			assert.JSONEq(t, cs.outJSON, w.Body.String())

			if cs.outStatus == http.StatusOK {
				var count int64
				time.Sleep(5 * time.Second) // maximum time of 4 seconds when inserting a record
				tdb.Model(&models.Callback{}).Count(&count)
				assert.NotZero(t, count, "At least one item should have been created")

				time.Sleep(5 * time.Second) // Give the test some extra time to remove the records
				tdb.Model(&models.Callback{}).Count(&count)
				assert.Zero(t, count, "Every item should have been removed after 7 seconds")
			}
		})
	}
}
