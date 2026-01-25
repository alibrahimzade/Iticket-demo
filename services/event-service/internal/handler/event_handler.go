package handler

import (
	"encoding/json"
	"net/http"

	"event-service/internal/dto"
	"event-service/internal/service"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}
func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events := h.service.GetAllEvents()

	response := make([]dto.EventResponse, 0, len(events))

	for _, e := range events {
		response = append(response, dto.EventResponse{
			ID:       e.ID,
			Title:    e.Title,
			Category: e.Category,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
