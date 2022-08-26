package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

func main() {
	//app := new(base.App)
	//app.Init()
	//
	//app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Hello World!")
	//}).Methods("GET")
	//
	//app.Run("8080")

	http.HandleFunc("/", echo)
	http.ListenAndServe(":8080", nil)
}

func echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь с клиентом прервана
		}

		connection.WriteMessage(websocket.TextMessage, message)

		go messageHandler(message)
	}
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}
