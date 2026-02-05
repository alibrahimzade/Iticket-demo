package handler

import (
	"encoding/json"
	"net/http"
	"ticket-service/internal/dto"
	"ticket-service/internal/service"
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
