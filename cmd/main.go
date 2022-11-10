package main

import (
	"fmt"
	"github.com/pavel-one/WebhookWatcher/internal/base"
	"github.com/pavel-one/WebhookWatcher/internal/controllers"
	"log"
	"os"
	"os/exec"
	"runtime"
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
	hunterController.Init(socketChannel)

	requestController := new(controllers.RequestController)
	requestController.Init(socketChannel)

	app.Router.Use(controllers.LoggingMiddleware)
	app.Router.NotFoundHandler = controllers.Handler404

	//api
	app.GET("/api/v1/", hunterController.Index)
	app.POST("/api/v1/", hunterController.Create)
	app.POST("/api/v1/check/", hunterController.Check)
	app.GET("/api/v1/channels/", hunterController.GetChannels)
	app.DELETE("/api/v1/channels/{id:[0-9]+}", hunterController.DropChannel)
	app.Prefix("/api/v1/request/", requestController.NewRequest)

	//load static files
	app.Static("/", "frontend")

	//websocket
	socket.Router.Use(controllers.LoggingMiddleware)
	socket.Router.NotFoundHandler = controllers.Handler404
	socket.GET("/", socketController.Connect)
	socket.GET("/{channel:[a-zA-Z0-9]+}", socketController.Connect)

	go app.ApiRun("3000", fatalChan)
	go socket.ApiRun("8080", fatalChan)

	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", "http://localhost:3000").Start()
		fmt.Println("Open http://localhost:3000 in browser...")
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:3000").Start()
		fmt.Println("Open http://localhost:3000 in browser...")
	case "darwin":
		exec.Command("open", "http://localhost:3000").Start()
		fmt.Println("Open http://localhost:3000 in browser...")
	}

	err := <-fatalChan
	if err != nil {
		app.Close()
		log.Printf("[FATAL] %v", err)
		os.Exit(1)
	}
}
