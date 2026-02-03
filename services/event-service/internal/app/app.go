package app

import (
	"context"
	"log"
	"net/http"
)

type App struct {
	server *http.Server
}

func New(server *http.Server) *App {
	return &App{server: server}
}
func (a *App) Run() {
	go func() {
		log.Println("HTTP server started")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Println("shutting down server...")
	return a.server.Shutdown(ctx)
}
