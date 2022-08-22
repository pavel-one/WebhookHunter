package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("[ERROR] failed to load environment: %v", err)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "Hello, world!")
		return
	})

	http.Handle("/", r)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error start server: %s", err.Error()))
		return
	}
}
