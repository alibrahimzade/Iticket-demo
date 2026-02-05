package service

import (
	"context"
	"ticket-service/internal/model"
	"ticket-service/internal/repository"

	"github.com/google/uuid"
)

type TicketService struct {
	repo repository.TicketRepository
}

func NewTicketService(repo repository.TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

func (s *TicketService) CreateTicket(
	ctx context.Context,
	eventID, zone string,
	price float64,
	currency string) (model.Ticket, error) {

	ticket := model.Ticket{
		ID:       uuid.NewString(),
		EventID:  eventID,
		Zone:     zone,
		Price:    price,
		Currency: currency,
		Status:   "AVAILABLE",
	}
	return s.repo.Create(ctx, ticket)
}

func (s *TicketService) GetTicketsByEvent(
	ctx context.Context, eventID string) ([]model.Ticket, error) {
	return s.repo.FindByEventID(ctx, eventID)
}

func (s *TicketService) UpdateTicket(
	ctx context.Context,
	id string,
	zone string,
	price float64,
	currency string,
	status string,
) (model.Ticket, error) {

	if id == "" || zone == "" || currency == "" || status == "" {
		return model.Ticket{}, ErrInvalidInput
	}

	return s.repo.Update(ctx, model.Ticket{
		ID:       id,
		Zone:     zone,
		Price:    price,
		Currency: currency,
		Status:   status,
	})
}

func (s *TicketService) ReserveTicket(
	ctx context.Context,
	id string,
) (model.Ticket, error) {
	return s.repo.Reserve(ctx, id)
}
