package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"event-service/internal/dto"
	"event-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()

	page, _ := strconv.Atoi(qp.Get("page"))
	size, _ := strconv.Atoi(qp.Get("size"))

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	category := qp.Get("category")

	// sorting
	sortBy := qp.Get("sort")
	if sortBy == "" {
		sortBy = "id"
	}

	order := qp.Get("order")
	if order == "" {
		order = "asc"
	}

	events, total, err := h.service.GetEvents(
		r.Context(),
		service.EventQuery{
			Category: category,
			Page:     page,
			Size:     size,
			SortBy:   sortBy,
			Order:    order,
		},
	)
	if err != nil {
		respondError(
			w,
			r,
			http.StatusInternalServerError,
			"EVENTS_FETCH_FAILED",
			err.Error(),
			nil,
		)
		return
	}

	items := make([]dto.EventResponse, 0, len(events))
	for _, e := range events {
		items = append(items, dto.EventResponse{
			ID:       e.ID,
			Title:    e.Title,
			Category: e.Category,
		})
	}

	respondJSON(w, http.StatusOK, dto.EventListResponse{
		Items: items,
		Page:  page,
		Size:  size,
		Total: total,
	})
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(
			w,
			r,
			http.StatusBadRequest,
			"INVALID_JSON",
			"invalid request body",
			nil,
		)
		return
	}

	created, err := h.service.CreateEvent(
		r.Context(),
		req.Title,
		req.Category,
	)
	if err != nil {
		respondError(
			w,
			r,
			http.StatusBadRequest,
			"EVENT_CREATE_FAILED",
			err.Error(),
			nil,
		)
		return
	}

	resp := dto.CreateEventResponse{
		ID:       created.ID,
		Title:    created.Title,
		Category: created.Category,
	}

	respondJSON(w, http.StatusCreated, resp)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, r, http.StatusBadRequest, "INVALID_ID", "event id is required", nil)
		return
	}

	var req dto.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "INVALID_JSON", "invalid request body", nil)
		return
	}

	updated, err := h.service.UpdateEvent(
		r.Context(),
		id,
		req.Title,
		req.Category,
	)
	if err != nil {
		status := http.StatusBadRequest
		if err == service.ErrNotFound {
			status = http.StatusNotFound
		}
		respondError(w, r, status, "EVENT_UPDATE_FAILED", err.Error(), nil)
		return
	}

	respondJSON(w, http.StatusOK, dto.EventResponse{
		ID:       updated.ID,
		Title:    updated.Title,
		Category: updated.Category,
	})
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, r, http.StatusBadRequest, "INVALID_ID", "event id is required", nil)
		return
	}

	err := h.service.DeleteEvent(r.Context(), id)
	if err != nil {
		switch err {
		case service.ErrNotFound:
			respondError(w, r, http.StatusNotFound, "EVENT_NOT_FOUND", err.Error(), nil)
		default:
			respondError(w, r, http.StatusInternalServerError, "EVENT_DELETE_FAILED", err.Error(), nil)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
