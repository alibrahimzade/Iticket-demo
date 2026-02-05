package handler

import (
	"encoding/json"
	"net/http"
	"ticket-service/internal/dto"
	"ticket-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TicketHandler struct {
	service *service.TicketService
}

func NewTicketHandler(service *service.TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "INVALID_JSON", "invalid body", nil)
		return
	}

	created, err := h.service.CreateTicket(
		r.Context(),
		req.EventID,
		req.Zone,
		req.Price,
		req.Currency)

	if err != nil {
		respondError(w, r, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), nil)
		return
	}
	respondJSON(w, http.StatusCreated, dto.TicketResponse{
		ID:       created.ID,
		EventID:  created.EventID,
		Zone:     created.Zone,
		Price:    created.Price,
		Currency: created.Currency,
		Status:   created.Status,
	})
}
func (h *TicketHandler) GetTickets(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()

	eventID := qp.Get("eventId")
	if eventID == "" {
		respondError(
			w,
			r,
			http.StatusBadRequest,
			"EVENT_ID_REQUIRED",
			"eventId query param is required",
			nil)
		return
	}

	tickets, err := h.service.GetTicketsByEvent(r.Context(), eventID)
	if err != nil {
		respondError(
			w,
			r,
			http.StatusInternalServerError,
			"TICKETS_FETCH_FAILED",
			err.Error(),
			nil,
		)
		return
	}
	items := make([]dto.TicketResponse, 0, len(tickets))
	for _, t := range tickets {
		items = append(items, dto.TicketResponse{
			ID:       t.ID,
			EventID:  t.EventID,
			Zone:     t.Zone,
			Price:    t.Price,
			Currency: t.Currency,
			Status:   t.Status,
		})
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
		"total": len(items),
	})
}

func (h *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		respondError(
			w,
			r,
			http.StatusBadRequest,
			"INVALID_UUID",
			"ticket id must be a valid UUID",
			nil,
		)
		return
	}
	if id == "" {
		respondError(w, r, http.StatusBadRequest, "INVALID_ID", "ticket id is required", nil)
		return
	}

	var req dto.UpdateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "INVALID_JSON", "invalid request body", nil)
		return
	}

	updated, err := h.service.UpdateTicket(
		r.Context(),
		id,
		req.Zone,
		req.Price,
		req.Currency,
		req.Status,
	)

	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidInput {
			status = http.StatusBadRequest
		}

		respondError(w, r, status, "TICKET_UPDATE_FAILED", err.Error(), nil)
		return
	}

	respondJSON(w, http.StatusOK, dto.TicketResponse{
		ID:       updated.ID,
		EventID:  updated.EventID,
		Zone:     updated.Zone,
		Price:    updated.Price,
		Currency: updated.Currency,
		Status:   updated.Status,
	})
}

func (h *TicketHandler) ReserveTicket(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	t, err := h.service.ReserveTicket(r.Context(), id)

	if err != nil {
		respondError(w, r, http.StatusConflict, "TICKET_NOT_AVAILABLE", err.Error(), nil)
		return
	}

	respondJSON(w, http.StatusOK, t)
}
