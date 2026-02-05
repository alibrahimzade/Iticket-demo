package repository

import (
	"context"
	"database/sql"
	"ticket-service/internal/model"
)

type PostgresTicketRepository struct {
	db *sql.DB
}

func NewPostgresTicketRepository(db *sql.DB) *PostgresTicketRepository {
	return &PostgresTicketRepository{db: db}
}

func (r *PostgresTicketRepository) Create(
	ctx context.Context, t model.Ticket,
) (model.Ticket, error) {

	query := `
		INSERT INTO tickets(id,event_id,zone,price,currency,status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, event_id, zone, price, currency,status, created_at
		`

	var created model.Ticket

	err := r.db.QueryRowContext(
		ctx,
		query,
		t.ID,
		t.EventID,
		t.Zone,
		t.Price,
		t.Currency,
		t.Status,
	).Scan(
		&created.ID,
		&created.EventID,
		&created.Zone,
		&created.Price,
		&created.Currency,
		&created.Status,
		&created.CreatedAt,
	)

	return created, err
}

func (r *PostgresTicketRepository) FindByEventID(
	ctx context.Context, eventID string) ([]model.Ticket, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, event_id, zone, price, currency, status, created_at
				FROM tickets WHERE event_id = $1`,
		eventID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		var t model.Ticket
		if err := rows.Scan(
			&t.ID,
			&t.EventID,
			&t.Zone,
			&t.Price,
			&t.Currency,
			&t.Status,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}

		tickets = append(tickets, t)
	}
	return tickets, nil
}
