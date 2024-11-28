package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/TOMMy-Net/effective-mobile/docs"
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
	defer file.Close()

	l := logger.InitLogger(file)

	service := handlers.NewService() // create service instance
	service.Storage = store
	service.Log = l

	// ---- routeres
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler) // http://host:port/swagger/index.html

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/songs", service.SongHandlers()).Methods(http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch)
	api.HandleFunc("/songs/{id}/text", service.GetSongTextHandler()).Methods(http.MethodGet)
	
	router.Use(middleware.ScanTrafic(l))
	// ---- end

	srv := newServer(router) // http server

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	wg.Add(1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
		log.Println("server Shutdown")
		wg.Done()
	}()

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()
	wg.Wait()
}

func newServer(r http.Handler) http.Server {
	if os.Getenv("PORT") == "" {
		log.Fatal("no server port")
	}

	return http.Server{
		Addr:         os.Getenv("PORT"),
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
