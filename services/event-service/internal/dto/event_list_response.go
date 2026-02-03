package dto

type EventListResponse struct {
	Items []EventResponse `json:"items"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
}
