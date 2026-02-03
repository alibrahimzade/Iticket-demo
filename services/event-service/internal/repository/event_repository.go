package repository

import (
	"context"
	"event-service/internal/model"
)

type EventRepository interface {
	FindAll(
		ctx context.Context,
		category string,
		sortBy string,
		order string,
		limit int,
		offset int,
	) ([]model.Event, error)

	Count(ctx context.Context, category string) (int, error)

	Create(
		ctx context.Context,
		event model.Event,
	) (model.Event, error)

	Update(ctx context.Context, event model.Event) (model.Event, error)
	Delete(ctx context.Context, id string) error
}
