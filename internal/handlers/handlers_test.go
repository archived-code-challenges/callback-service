package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/noelruault/go-callback-service/internal/models"
	"github.com/noelruault/go-callback-service/internal/web"
)

type testCallbackService struct {
	models.CallbackService
	upsert func(context.Context, []models.Callback) error
}

func (t *testCallbackService) Upsert(ctx context.Context, cs []models.Callback) error {
	if t.upsert != nil {
		return t.upsert(ctx, cs)
	}

	panic("not provided")
}

func NewTestContext() context.Context {
	return context.WithValue(context.Background(), web.KeyValues, &web.Values{})
}

func TestCallback_Handle(t *testing.T) {
	csvc := &testCallbackService{}
	c := NewCallbacks(csvc, nil)

	var cases = []struct {
		name      string
		input     string
		outStatus int
		outJSON   string
		setup     func(*testing.T)
	}{
		{
			"notJSON",
			"a dalhd lkald fkjahd lfkjasdlf ",
			http.StatusBadRequest,
			`{"error":"invalid_json","message":"provided input cannot be parsed"}`,
			nil,
		},
		{
			"removeDuplicateValues",
			`{"object_ids": [91,10,78,91,10,78]}`,
			http.StatusOK,
			`{}`,
			func(t *testing.T) {
				csvc.upsert = func(ctx context.Context, cs []models.Callback) error {
					assert.Len(t, cs, 3)
					return nil
				}
			},
		},
		{
			"ok",
			`{"object_ids": [91,10,78]}`,
			http.StatusOK,
			`{}`,
			func(t *testing.T) {
				csvc.upsert = func(ctx context.Context, cs []models.Callback) error {
					assert.Len(t, cs, 3)
					return nil
				}
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewReader([]byte(cs.input)))
			ctx := NewTestContext()

			if cs.setup != nil {
				cs.setup(t)
			}

			c.Handle(ctx, w, r)

			assert.Equal(t, cs.outStatus, w.Result().StatusCode)
			assert.JSONEq(t, cs.outJSON, w.Body.String())

			*csvc = testCallbackService{}
		})
	}
}
