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
	"time"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	log.Println("stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error(err)
	}

	log.Println("closing database...")
	db.Close()
}
