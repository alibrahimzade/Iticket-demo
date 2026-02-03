package dto

type CreateEventResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
}
