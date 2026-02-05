package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"event-service/internal/app"
	"event-service/internal/config"
	"event-service/internal/handler"
	"event-service/internal/repository"
	"event-service/internal/service"
)

// @title           Event Service API
// @version         1.0
// @description     Event service for Iticket platform
// @termsOfService  http://example.com/terms/

// @contact.name   Backend Team
// @contact.email  backend@iticket.io

// @host      localhost:8081
// @BasePath  /api/v1
func main() {

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
	slog.SetDefault(logger)

	// 1️⃣ Load config
	cfg := config.Load()

	if cfg.DBURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// 2️⃣ Initialize Postgres
	db, err := repository.NewPostgresDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	app.RunMigrations(db, "file://migrations")
	
	defer func() {
		slog.Info("closing database connection")
		if err := db.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	// 3️⃣ Build dependencies
	eventRepo := repository.NewPostgresEventRepository(db)
	eventService := service.NewEventService(eventRepo)
	eventHandler := handler.NewEventHandler(eventService)
	readinessHandler := handler.NewReadinessHandler(db)

	// 4️⃣ Setup router
	router := setupRouter(eventHandler, readinessHandler)

	log.Printf("Starting server on address = [%q]", ":"+cfg.Port)

	server := &http.Server{
		Addr:    ":" + cfg.Port, // MUST start with colon
		Handler: router,
	}

	application := app.New(server)
	application.Run()

	slog.Info("Event service started", "port", cfg.Port)

	// 6️⃣ Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("Event Service stopped")
}

func setupRouter(eventHandler *handler.EventHandler,
	readinessHandler *handler.ReadinessHandler) http.Handler {

	r := chi.NewRouter()

	// middleware
	r.Use(handler.RequestIDMiddleware)
	r.Use(handler.RecoveryMiddleware)
	r.Use(handler.LoggingMiddleware)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Get("/ready", readinessHandler.Ready)

	// api
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/events", eventHandler.GetEvents)
		r.Post("/events", eventHandler.CreateEvent)
		r.Put("/events/{id}", eventHandler.UpdateEvent)
		r.Delete("/events/{id}", eventHandler.DeleteEvent)
	})

	return r
}
