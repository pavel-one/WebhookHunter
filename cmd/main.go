package main

import (
	"github.com/pavel-one/WebhookWatcher/internal/base"
	"github.com/pavel-one/WebhookWatcher/internal/controllers"
	"log"
	"os"
)

func main() {
	fatalChan := make(chan error, 1)

	app := new(base.App)
	app.Init()

	hunterController := new(controllers.HunterController)
	hunterController.Init(app.DB)

	app.Router.Use(controllers.LoggingMiddleware)
	app.Router.NotFoundHandler = controllers.Handler404
	app.GET("/", hunterController.Index)
	app.POST("/", hunterController.Create)

	go app.ApiRun("8080", fatalChan)

	err := <-fatalChan
	if err != nil {
		app.Close()
		log.Printf("[FATAL] %v", err)
		os.Exit(1)
	}
}
