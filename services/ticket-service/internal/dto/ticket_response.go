package dto

type TicketResponse struct {
	ID       string  `json:"id"`
	EventID  string  `json:"event_id"`
	Zone     string  `json:"zone"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}
