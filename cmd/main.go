package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pavel-one/WebhookWatcher/internal/base"
	"github.com/pavel-one/WebhookWatcher/internal/controllers"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
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

	// all routes without /api /ui
	app.Router.MatcherFunc(func(request *http.Request, match *mux.RouteMatch) bool {
		find, _ := regexp.MatchString(`^/(\bapi|ui\b.*).*$`, request.RequestURI)
		if find {
			return false
		}

		return true
	}).PathPrefix("/").HandlerFunc(requestController.NewRequest)

	//api
	app.GET("/api/v1/", hunterController.Index)
	app.POST("/api/v1/", hunterController.Create)
	app.POST("/api/v1/check/", hunterController.Check)
	app.GET("/api/v1/channels/", hunterController.GetChannels)
	app.DELETE("/api/v1/channels/{id:[0-9]+}", hunterController.DropChannel)

	//load static files
	app.Static("/ui", "frontend")

	//websocket
	socket.Router.Use(controllers.LoggingMiddleware)
	socket.Router.NotFoundHandler = controllers.Handler404
	socket.Prefix("/", socketController.Connect)

	go app.Run("3000", fatalChan)
	go socket.Run("8080", fatalChan)

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
