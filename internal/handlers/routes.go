package handlers

import (
	"log"
	"net/http"

	"gorm.io/gorm"

	mw "github.com/noelruault/go-callback-service/internal/middleware"
	"github.com/noelruault/go-callback-service/internal/models"
	"github.com/noelruault/go-callback-service/internal/web"
)

func API(log *log.Logger, db *gorm.DB, callbackServiceURL string) http.Handler {
	app := web.NewApp(log, mw.Logger(log), mw.Metrics(), mw.Panics(log))

	// Models
	cm := models.NewCallbackService(db, callbackServiceURL, log)

	{
		c := Check{db: db, log: log}
		app.Handle(http.MethodGet, "/", c.Health)
	}
	// Handlers
	{
		csvc := NewCallbacks(cm, log)
		app.Handle(http.MethodPost, "/callback", csvc.Handle)
	}

	return app
}
