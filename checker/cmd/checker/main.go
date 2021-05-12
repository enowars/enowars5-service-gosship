package main

import (
	"checker/pkg/checker"
	"checker/pkg/handler"
	"context"
	"gosship/pkg/logger"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	log := logger.New()
	checkerHandler := handler.New()
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: checker.NewChecker(log, checkerHandler),
	}
	go func() {
		log.Printf("starting server on port %s...", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}
	log.Println("server stopped")
}
