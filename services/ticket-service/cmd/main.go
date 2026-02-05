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

	"ticket-service/internal/app"
	"ticket-service/internal/config"
	"ticket-service/internal/handler"
	"ticket-service/internal/repository"
	"ticket-service/internal/service"
)

// @title           Ticket Service API
// @version         1.0
// @description     Ticket service for Iticket platform
// @termsOfService  http://example.com/terms/

// @contact.name   Backend Team
// @contact.email  backend@iticket.io

// @host      localhost:8082
// @BasePath  /api/v1
func main() {

	// ---------- Logger ----------
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
	slog.SetDefault(logger)

	// ---------- Load config ----------
	cfg := config.Load()

	if cfg.DBURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// ---------- Initialize Postgres ----------
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

	// ---------- Build dependencies (Ticket domain) ----------
	ticketRepo := repository.NewPostgresTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketService)

	readinessHandler := handler.NewReadinessHandler(db)

	// ---------- Setup router ----------
	router := setupRouter(ticketHandler, readinessHandler)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	application := app.New(server)
	application.Run()

	slog.Info("Ticket service started", "port", cfg.Port)

	// ---------- Graceful shutdown ----------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("Ticket Service stopped")
}

func setupRouter(
	ticketHandler *handler.TicketHandler,
	readinessHandler *handler.ReadinessHandler,
) http.Handler {

	r := chi.NewRouter()

	// ---------- middleware ----------
	r.Use(handler.RequestIDMiddleware)
	r.Use(handler.RecoveryMiddleware)
	r.Use(handler.LoggingMiddleware)

	// ---------- health ----------
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Get("/ready", readinessHandler.Ready)

	// ---------- API ----------
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/tickets", ticketHandler.CreateTicket)
		// GET /tickets?eventId=...  (we’ll add next)
		r.Get("/tickets", ticketHandler.GetTickets)
	})

	return r
}
