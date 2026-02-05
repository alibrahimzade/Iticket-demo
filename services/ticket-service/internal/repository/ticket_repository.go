package repository

import (
	"context"
	"ticket-service/internal/model"
)

type TicketRepository interface {
	Create(ctx context.Context, ticket model.Ticket) (model.Ticket, error)
	FindByEventID(ctx context.Context, eventID string) ([]model.Ticket, error)
	Update(ctx context.Context, ticket model.Ticket) (model.Ticket, error)
	Reserve(ctx context.Context, id string) (model.Ticket, error)
}
