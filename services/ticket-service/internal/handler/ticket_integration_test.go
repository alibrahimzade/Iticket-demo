package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"ticket-service/internal/handler"
	"ticket-service/internal/repository"
	"ticket-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func setupTicketTestServer(t *testing.T) http.Handler {
	db, err := repository.NewPostgresDB(
		"postgres://iticket:iticket@localhost:5432/iticket?sslmode=disable",
	)

	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM tickets")
		_ = db.Close()
	})

	repo := repository.NewPostgresTicketRepository(db)
	svc := service.NewTicketService(repo)
	h := handler.NewTicketHandler(svc)

	r := chi.NewRouter()
	r.Post("/api/v1/tickets", h.CreateTicket)
	r.Get("/api/v1/tickets", h.GetTickets)
	r.Put("/api/v1/tickets/{id}", h.UpdateTicket)

	return r
}

func TestCreateAndGetTickets(t *testing.T) {
	server := setupTicketTestServer(t)

	createReq := map[string]interface{}{
		"event_id": "50142040-e12c-4d9f-bd62-54a24e7fe841",
		"zone":     "VIP",
		"price":    50,
		"currency": "AZN",
	}

	body, _ := json.Marshal(createReq)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tickets",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	getReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tickets?eventId=50142040-e12c-4d9f-bd62-54a24e7fe841",
		nil,
	)

	getRec := httptest.NewRecorder()
	server.ServeHTTP(getRec, getReq)

	require.Equal(t, http.StatusOK, getRec.Code)
}

func TestUpdateTicket(t *testing.T) {
	server := setupTicketTestServer(t)

	create := []byte(`{
		"event_id":"50142040-e12c-4d9f-bd62-54a24e7fe841",
		"zone":"FAN",
		"price":30,
		"currency":"USD"
	}`)

	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(create))
	postReq.Header.Set("Content-Type", "application/json")

	postRec := httptest.NewRecorder()
	server.ServeHTTP(postRec, postReq)

	var created struct {
		ID string `json:"id"`
	}

	_ = json.NewDecoder(postRec.Body).Decode(&created)

	update := []byte(`{
		"zone":"VIP",
		"price":100,
		"currency":"USD",
		"status":"AVAILABLE"
	}`)

	putReq := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/tickets/"+created.ID,
		bytes.NewReader(update),
	)

	putReq.Header.Set("Content-Type", "application/json")

	putRec := httptest.NewRecorder()
	server.ServeHTTP(putRec, putReq)

	require.Equal(t, http.StatusOK, putRec.Code)
}
