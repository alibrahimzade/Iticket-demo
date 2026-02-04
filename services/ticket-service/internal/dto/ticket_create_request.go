package dto

type CreateTicketRequest struct {
	EventID  string  `json:"event_id"`
	Zone     string  `json:"zone"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}
