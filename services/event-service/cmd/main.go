package main

import (
	"event-service/internal/handler"
	"event-service/internal/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	eventService := service.NewEventService()
	eventHandler := handler.NewEventHandler(eventService)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/events", eventHandler.GetEvents)

	log.Println("Event Service started on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
