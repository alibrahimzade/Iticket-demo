package dto

type CreateEventRequest struct {
	Title    string `json:"title"`
	Category string `json:"category"`
}
