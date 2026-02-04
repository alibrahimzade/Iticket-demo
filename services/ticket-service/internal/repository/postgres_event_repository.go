package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"ticket-service/internal/model"
)

type PostgresEventRepository struct {
	db *sql.DB
}

func NewPostgresEventRepository(db *sql.DB) *PostgresEventRepository {
	return &PostgresEventRepository{db: db}
}

func (r *PostgresEventRepository) FindAll(
	ctx context.Context,
	category string,
	sortBy string,
	order string,
	limit int,
	offset int,
) ([]model.Event, error) {

	// 🔐 whitelist columns
	sortColumn := "id"
	switch sortBy {
	case "title":
		sortColumn = "title"
	case "category":
		sortColumn = "category"
	}

	sortOrder := "ASC"
	if strings.ToLower(order) == "desc" {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf(`
		SELECT id, title, category
		FROM events
		WHERE ($1 = '' OR category = $1)
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortColumn, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Category); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, rows.Err()
}

func (r *PostgresEventRepository) Count(
	ctx context.Context,
	category string,
) (int, error) {

	query := `
		SELECT COUNT(*)
		FROM events
		WHERE ($1 = '' OR category = $1)
	`

	var total int
	err := r.db.QueryRowContext(ctx, query, category).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *PostgresEventRepository) Create(
	ctx context.Context,
	event model.Event,
) (model.Event, error) {

	query := `
		INSERT INTO events (id, title, category)
		VALUES ($1, $2, $3)
		RETURNING id, title, category
	`

	var created model.Event
	err := r.db.QueryRowContext(
		ctx,
		query,
		event.ID,
		event.Title,
		event.Category,
	).Scan(&created.ID, &created.Title, &created.Category)

	if err != nil {
		return model.Event{}, err
	}

	return created, nil
}

func (r *PostgresEventRepository) Update(
	ctx context.Context,
	event model.Event,
) (model.Event, error) {

	query := `
		UPDATE events
		SET title = $2, category = $3
		WHERE id = $1
		RETURNING id, title, category
	`

	var updated model.Event
	err := r.db.QueryRowContext(
		ctx,
		query,
		event.ID,
		event.Title,
		event.Category,
	).Scan(&updated.ID, &updated.Title, &updated.Category)

	if err == sql.ErrNoRows {
		return model.Event{}, ErrNotFound
	}
	if err != nil {
		return model.Event{}, err
	}

	return updated, nil
}

func (r *PostgresEventRepository) Delete(
	ctx context.Context,
	id string,
) error {

	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM events WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
