package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCallbackGorm_Upsert(t *testing.T) {
	cdb := callbackGorm{NewTestDatabase(t)}
	CallbackSelfDeleteTime = 800 * time.Millisecond
	defer func() {
		CallbackSelfDeleteTime = 30 * time.Second
		CleanupTestDatabase(cdb.db)
	}()

	var cases = []struct {
		name         string
		callbacks    []Callback
		outcallbacks []Callback
		outerr       error
		setup        func(*testing.T)
	}{
		{
			"ok",
			[]Callback{{ID: 123, Online: true, Timestamp: 1111111}},
			[]Callback{{ID: 123, Online: true, Timestamp: 1111111}},
			nil,
			nil,
		},
		{
			"IDisDuplicatedValue",
			[]Callback{{ID: 123, Online: true, Timestamp: 1111111}},
			[]Callback{{ID: 123, Online: true, Timestamp: 1111111}},
			nil,
			func(t *testing.T) {
				cdb.db.Create(&Callback{ID: 123, Online: true, Timestamp: 1111111})
			},
		},
		{
			"updatedWhenConflict",
			[]Callback{{ID: 123, Online: true, Timestamp: 1111111}},
			[]Callback{{ID: 123, Online: true, Timestamp: 1111111}},
			nil,
			func(t *testing.T) {
				cdb.db.Create(&Callback{ID: 123, Online: false, Timestamp: 0})
			},
		},
	}
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			ctx := context.Background()

			if cs.setup != nil {
				cs.setup(t)
			}

			cdb.Upsert(ctx, cs.callbacks)

			// Check if the callback could have been created
			var callbacksDB []Callback
			cdb.db.Find(&callbacksDB)
			assert.Equal(t, cs.outcallbacks, callbacksDB)

			// Check if the entity is removed after self-delete time
			var count int64
			time.Sleep(1500 * time.Millisecond)
			cdb.db.Model(&Callback{}).Count(&count)
			assert.Zero(t, count)

			CleanupTestDatabase(cdb.db)
		})
	}
}
