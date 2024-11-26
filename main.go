package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TOMMy-Net/effective-mobile/internal/handlers"
	"github.com/TOMMy-Net/effective-mobile/internal/middleware"
	"github.com/TOMMy-Net/effective-mobile/internal/storage/db"
	"github.com/TOMMy-Net/effective-mobile/tools"
	"github.com/TOMMy-Net/effective-mobile/tools/logger"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

)

// @title Music Library API
// @version 1.0
// @description API for managing music library

// @host localhost:8000
// @BasePath /
func main() {
	tools.NewValidator() // init singletone validator
	if err := tools.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	store, err := db.ConnectPostgres(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.DB.Close()

	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	l := logger.InitLogger(file)

	service := handlers.NewService() // create service instance
	service.Storage = store
	service.Log = l

	router := mux.NewRouter()
	router.HandleFunc("/songs", service.SongHandlers()).Methods(http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch)
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	router.Use(middleware.ScanTrafic(l))

	srv := http.Server{
		Addr:         ":8000",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}
