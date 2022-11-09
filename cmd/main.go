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

	socket := new(base.App)
	socket.Init()

	socketChannel := make(chan controllers.SocketMessage, 10)
	socketController := new(controllers.SocketController)
	socketController.Init(socketChannel)

	hunterController := new(controllers.HunterController)
	hunterController.Init()

	requestController := new(controllers.RequestController)
	requestController.Init(socketChannel)

	app.Router.Use(controllers.LoggingMiddleware)
	app.Router.NotFoundHandler = controllers.Handler404

	//load static files
	app.Static("/web/", "frontend")

	//api
	app.GET("/", hunterController.Index)
	app.POST("/", hunterController.Create)
	app.POST("/check/", hunterController.Check)
	app.GET("/channels/", hunterController.GetChannels)
	app.Prefix("/request/", requestController.NewRequest)

	//websocket
	socket.Router.Use(controllers.LoggingMiddleware)
	socket.Router.NotFoundHandler = controllers.Handler404
	socket.GET("/", socketController.Connect)
	socket.GET("/{channel:[a-zA-Z0-9]+}", socketController.Connect)

	go app.ApiRun("80", fatalChan)
	go socket.ApiRun("8080", fatalChan)

	err := <-fatalChan
	if err != nil {
		app.Close()
		log.Printf("[FATAL] %v", err)
		os.Exit(1)
	}
}
