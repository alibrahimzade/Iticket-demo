package service

import (
	"context"

	"event-service/internal/model"
	"event-service/internal/repository"

	"github.com/google/uuid"
)

type EventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetEvents(
	ctx context.Context,
	q EventQuery,
) ([]model.Event, int, error) {

	offset := (q.Page - 1) * q.Size

	events, err := s.repo.FindAll(
		ctx,
		q.Category,
		q.SortBy,
		q.Order,
		q.Size,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, q.Category)
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (s *EventService) CreateEvent(
	ctx context.Context,
	title string,
	category string,
) (model.Event, error) {

	if title == "" {
		return model.Event{}, ErrInvalidInput
	}
	if category == "" {
		return model.Event{}, ErrInvalidInput
	}

	event := model.Event{
		ID:       uuid.NewString(),
		Title:    title,
		Category: category,
	}

	return s.repo.Create(ctx, event)
}

func (s *EventService) UpdateEvent(
	ctx context.Context,
	id string,
	title string,
	category string,
) (model.Event, error) {

	if title == "" || category == "" {
		return model.Event{}, ErrInvalidInput
	}

	return s.repo.Update(ctx, model.Event{
		ID:       id,
		Title:    title,
		Category: category,
	})
}

func (s *EventService) DeleteEvent(
	ctx context.Context,
	id string,
) error {
	return s.repo.Delete(ctx, id)
}
