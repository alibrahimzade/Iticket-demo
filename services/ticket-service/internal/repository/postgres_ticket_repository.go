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

func (r *PostgresTicketRepository) Update(
	ctx context.Context,
	t model.Ticket,
) (model.Ticket, error) {

	query := `
		UPDATE tickets
		SET zone = $2,
			price = $3,
			currency = $4,
			status = $5
		WHERE id = $1
		RETURNING id, event_id, zone, price, currency, status, created_at
		`

	var updated model.Ticket
	err := r.db.QueryRowContext(
		ctx,
		query,
		t.ID,
		t.Zone,
		t.Price,
		t.Status,
		t.Currency,
	).Scan(
		&updated.ID,
		&updated.EventID,
		&updated.Zone,
		&updated.Price,
		&updated.Currency,
		&updated.Status,
		&updated.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return model.Ticket{}, ErrNotFound
	}

	if err != nil {
		return model.Ticket{}, err
	}

	return updated, nil
}

func (r *PostgresTicketRepository) Reserve(
	ctx context.Context,
	id string,
) (model.Ticket, error) {
	query := `
		UPDATE tickets
		SET status = 'RESERVED'
		WHERE id = $1 AND status = 'AVAILABLE'
		RETURNING id, event_id, zone, price, currency, status`

	var t model.Ticket
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.EventID, &t.Zone, &t.Price, &t.Currency, &t.Status)

	if err == sql.ErrNoRows {
		return model.Ticket{}, ErrInvalidState
	}

	if err != nil {
		return model.Ticket{}, err
	}

	return t, nil
}
