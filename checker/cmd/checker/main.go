package main

import (
	"checker/pkg/checker"
	"checker/pkg/database"
	"checker/pkg/handler"
	"context"
	"gosship/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := logger.New()
	db, err := database.NewDatabase(log)
	if err != nil {
		log.Fatal(err)
	}

	checkerHandler := handler.New(log, db)
	server := &http.Server{
		Addr:    ":8080",
		Handler: checker.NewChecker(log, checkerHandler),
	}
	go func() {
		log.Printf("starting server on port %s...", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("stopping server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}
	log.Println("closing database...")
	db.Close()
}
