package dto

type UpdateTicketRequest struct {
	Zone     string  `json:"zone"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}
