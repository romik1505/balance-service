package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/romik1505/balance-service/internal/config"
)

type App struct {
	httpServer http.Server
}

func NewApp(ctx context.Context, handler http.Handler) *App {
	port := config.GetValue(config.Port)
	return &App{
		httpServer: http.Server{
			Addr:         port,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			Handler:      handler,
		},
	}
}

func (a *App) Run() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	go func(ch chan os.Signal) {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Println(err.Error())
			done <- os.Interrupt
			return
		}
	}(done)

	log.Printf("Server started on %s port", config.GetValue(config.Port))

	<-done
	defer close(done)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	log.Println("Server gracefully closed")
	return a.httpServer.Shutdown(ctx)
}
