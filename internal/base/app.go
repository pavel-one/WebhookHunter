package base

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
	DB     *sqlx.DB
	Router *mux.Router
}

func (a *App) Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("[FATAL] Not loading environment: %v", err)
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	))

	if err != nil {
		log.Fatalf("[FATAL] Unable to connect to database: %v", err)
	}

	a.DB = db
	a.Router = mux.NewRouter()
}

func (a *App) Run(port string) {
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("[DEBUG] Running server on port %s", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(fmt.Sprintf("[FATAL] Error start server: %s", err.Error()))
	}
}
