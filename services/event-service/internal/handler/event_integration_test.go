package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"event-service/internal/handler"
	"event-service/internal/repository"
	"event-service/internal/service"
)

func setUpTestServer(t *testing.T) http.Handler {
	db, err := repository.NewPostgresDB(
		"postgres://iticket:iticket@localhost:5432/iticket?sslmode=disable",
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM events")
		_ = db.Close()
	})

	repo := repository.NewPostgresEventRepository(db)
	svc := service.NewEventService(repo)
	h := handler.NewEventHandler(svc)

	r := chi.NewRouter()
	r.Post("/api/v1/events", h.CreateEvent)
	r.Get("/api/v1/events", h.GetEvents)
	r.Put("/api/v1/events/{id}", h.UpdateEvent)
	r.Delete("/api/v1/events/{id}", h.DeleteEvent)

	return r
}

func TestCreateAndGetEvents(t *testing.T) {
	server := setUpTestServer(t)

	reqBody := map[string]string{
		"title":    "Go Conference",
		"category": "tech",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/events",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var createResp struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Category string `json:"category"`
	}

	err = json.NewDecoder(rec.Body).Decode(&createResp)
	require.NoError(t, err)
	require.NotEmpty(t, createResp.ID)
	require.Equal(t, "Go Conference", createResp.Title)

	// ----------------GEt_--------
	getReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/events?page=1&size=10",
		nil)

	getRec := httptest.NewRecorder()
	server.ServeHTTP(getRec, getReq)

	require.Equal(t, http.StatusOK, getRec.Code)

	var listResp struct {
		Items []struct {
			ID       string `json:"id"`
			Title    string `json:"title"`
			Category string `json:"category"`
		} `json:"items"`
		Total int `json:"total"`
	}

	err = json.NewDecoder(getRec.Body).Decode(&listResp)
	require.NoError(t, err)

	require.Equal(t, 1, listResp.Total)
	require.Len(t, listResp.Items, 1)
	require.Equal(t, createResp.ID, listResp.Items[0].ID)
}

func TestUpdateEvent(t *testing.T) {
	server := setUpTestServer(t)

	//____-------Create-------_____

	createBody := []byte(`{"title":"Old Title","category":"old"}`)
	postReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/events",
		bytes.NewReader(createBody),
	)

	postReq.Header.Set("Content-Type", "application/json")

	postRec := httptest.NewRecorder()
	server.ServeHTTP(postRec, postReq)

	require.Equal(t, http.StatusCreated, postRec.Code)

	var created struct {
		ID string `json:"ID"`
	}
	require.NoError(t, json.NewDecoder(postRec.Body).Decode(&created))
	require.NotEmpty(t, created.ID)

	//-------------Updatre--------------
	updateBody := []byte(`{"title":"New Title","category":"new"}`)
	putReq := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/events/"+created.ID,
		bytes.NewReader(updateBody),
	)

	putReq.Header.Set("Content-Type", "application/json")

	putRec := httptest.NewRecorder()
	server.ServeHTTP(putRec, putReq)

	require.Equal(t, http.StatusOK, putRec.Code)

	var updated struct {
		Title    string `json:"title"`
		Category string `json:"category"`
	}
	require.NoError(t, json.NewDecoder(putRec.Body).Decode(&updated))

	require.Equal(t, "New Title", updated.Title)
	require.Equal(t, "new", updated.Category)
}

func TestDelete(t *testing.T) {
	server := setUpTestServer(t)

	body := []byte(`{"title":"Delete Me","category":"test"}`)
	postReq := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/events",
		bytes.NewReader(body),
	)

	postReq.Header.Set("Content-Type", "application/json")

	postRec := httptest.NewRecorder()
	server.ServeHTTP(postRec, postReq)

	require.Equal(t, http.StatusCreated, postRec.Code)

	var created struct {
		ID string `json:"id"`
	}

	require.NoError(t, json.NewDecoder(postRec.Body).Decode(&created))
	require.NotEmpty(t, created.ID)

	// ---- DELETE ----
	delReq := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/events/"+created.ID,
		nil,
	)

	delRec := httptest.NewRecorder()
	server.ServeHTTP(delRec, delReq)

	require.Equal(t, http.StatusNoContent, delRec.Code)

	// __--------Verify empty Get

	getReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/events?page=1&size=10",
		nil)

	getRec := httptest.NewRecorder()
	server.ServeHTTP(getRec, getReq)

	require.Equal(t, http.StatusOK, getRec.Code)

	var listResp struct {
		Total int `json:"total"`
	}

	require.NoError(t, json.NewDecoder(getRec.Body).Decode(&listResp))
	require.Equal(t, 0, listResp.Total)
}
