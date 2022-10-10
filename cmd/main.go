package main

import (
	"fmt"
	"github.com/pavel-one/WebhookWatcher/internal/base"
	"github.com/pavel-one/WebhookWatcher/internal/controllers"
	"log"
	"os"
	"time"
)

func main() {
	fatalChan := make(chan error, 1)

	app := new(base.App)
	app.Init()

	socket := new(base.App)
	socket.Init()

	socketController := new(controllers.SocketController)
	socketController.Init(socket.DB)

	hunterController := new(controllers.HunterController)
	hunterController.Init(app.DB)

	app.Router.Use(controllers.LoggingMiddleware)
	app.Router.NotFoundHandler = controllers.Handler404
	app.GET("/", hunterController.Index)
	app.POST("/check/", hunterController.Check)
	app.POST("/", hunterController.Create)

	socket.Router.Use(controllers.LoggingMiddleware)
	socket.Router.NotFoundHandler = controllers.Handler404
	socket.GET("/", socketController.Connect)
	socket.GET("/{channel:[a-Z]+}", socketController.Connect)

	go app.ApiRun("80", fatalChan)
	go socket.ApiRun("8080", fatalChan)

	// TODO: This example send group message, drop it
	go func() {
		i := 0

		for {
			socketController.MessageChain <- controllers.SocketMessage{
				Channel: "test",
				Message: fmt.Sprintf("Send test message # %d", i),
			}
			time.Sleep(time.Second * 2)
			i++
		}
	}()

	err := <-fatalChan
	if err != nil {
		app.Close()
		log.Printf("[FATAL] %v", err)
		os.Exit(1)
	}
}
