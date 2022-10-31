package main

import (
	"fmt"
	"github.com/pavel-one/WebhookWatcher/internal/adminApi"
	"github.com/pavel-one/WebhookWatcher/internal/base"
	"github.com/pavel-one/WebhookWatcher/internal/controllers"
	"github.com/pavel-one/WebhookWatcher/internal/middlewars"
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

	requestController := new(controllers.RequestController)
	requestController.Init(app.DB)

	adminController := new(adminApi.AdminController)
	adminController.Init(app.DB)

	app.Router.Use(middlewars.LoggingMiddleware)
	app.Router.NotFoundHandler = middlewars.Handler404

	//load static files
	app.Static("/web/", "frontend")

	//api
	app.GET("/", hunterController.Index)
	app.POST("/", hunterController.Create)
	app.POST("/check/", hunterController.Check)
	app.POST("/request/{channel:[a-zA-Z0-9]+}", requestController.NewRequest)
	app.POST("/admin/login/", adminController.Login)

	//here must be admin routes
	app.AdminRouteWithMiddleware("/test/", "GET", adminController.Test)
	////hunter
	adminHunterController := new(adminApi.AdminHunterController)
	adminHunterController.Init(app.DB)
	app.AdminRouteWithMiddleware("/hunter/create/", "POST", adminHunterController.Create)
	app.AdminRouteWithMiddleware("/hunter/{slug}/", "GET", adminHunterController.Get)
	app.AdminRouteWithMiddleware("/hunters/", "GET", adminHunterController.GetAll)
	app.AdminRouteWithMiddleware("/hunter/{slug}/update/", "PATCH", adminHunterController.Update)
	app.AdminRouteWithMiddleware("/hunter/{slug}/delete/", "DELETE", adminHunterController.Delete)
	////channel
	adminChannelController := new(adminApi.AdminChannelController)
	adminChannelController.Init(app.DB)
	app.AdminRouteWithMiddleware("/channel/{slug}/create", "POST", adminChannelController.Create)
	app.AdminRouteWithMiddleware("/channel/{slug}/{channel}/", "GET", adminChannelController.Get)
	app.AdminRouteWithMiddleware("/channels/{slug}/", "GET", adminChannelController.GetAllByHunter)
	app.AdminRouteWithMiddleware("/channels/", "GET", adminChannelController.GetAll)
	app.AdminRouteWithMiddleware("/channel/{slug}/{channel}/", "PATCH", adminChannelController.Update)
	app.AdminRouteWithMiddleware("/channel/{slug}/{channel}/", "DELETE", adminChannelController.Delete)

	socket.Router.Use(middlewars.LoggingMiddleware)
	socket.Router.NotFoundHandler = middlewars.Handler404
	app.GET("/request/{channel:[a-zA-Z0-9]+}", requestController.NewRequest)

	//websocket
	socket.Router.Use(middlewars.LoggingMiddleware)
	socket.Router.NotFoundHandler = middlewars.Handler404
	socket.GET("/", socketController.Connect)
	socket.GET("/{channel:[a-zA-Z0-9]+}", socketController.Connect)

	go app.ApiRun("80", fatalChan)
	go socket.ApiRun("8080", fatalChan)

	// TODO: This example send group message, drop it
	go func() {
		i := 0

		for {
			socketController.MessageChain <- controllers.SocketMessage{
				Domain:  "test",
				Channel: "/",
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
