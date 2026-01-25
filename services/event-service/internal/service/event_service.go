package service

import "event-service/internal/model"

type EventService struct {
}

func NewEventService() *EventService {
	return &EventService{}
}

func (s *EventService) GetAllEvents() []model.Event {
	return []model.Event{
		{
			ID:       "1",
			Title:    "Rock Concert",
			Category: "Concert",
		},
		{
			ID:       "2",
			Title:    "Cinema Night",
			Category: "Film",
		},
	}
}
