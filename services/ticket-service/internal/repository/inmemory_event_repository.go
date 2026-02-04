package repository

import (
	"context"
	"sort"
	"strings"
	"ticket-service/internal/model"
)

type InMemoryEventRepository struct{}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{}
}

func (r *InMemoryEventRepository) FindAll(
	ctx context.Context,
	category string,
	sortBy string,
	order string,
	limit int,
	offset int,
) ([]model.Event, error) {

	all := []model.Event{
		{ID: "1", Title: "Rock Concert", Category: "concert"},
		{ID: "2", Title: "Cinema Night", Category: "film"},
		{ID: "3", Title: "Business Seminar", Category: "seminar"},
	}

	// filter
	filtered := make([]model.Event, 0)
	for _, e := range all {
		if category == "" || strings.EqualFold(e.Category, category) {
			filtered = append(filtered, e)
		}
	}

	// sort
	sort.Slice(filtered, func(i, j int) bool {
		switch sortBy {
		case "title":
			if order == "desc" {
				return filtered[i].Title > filtered[j].Title
			}
			return filtered[i].Title < filtered[j].Title
		case "category":
			if order == "desc" {
				return filtered[i].Category > filtered[j].Category
			}
			return filtered[i].Category < filtered[j].Category
		default:
			if order == "desc" {
				return filtered[i].ID > filtered[j].ID
			}
			return filtered[i].ID < filtered[j].ID
		}
	})

	// paginate
	if offset >= len(filtered) {
		return []model.Event{}, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

func (r *InMemoryEventRepository) Count(
	ctx context.Context,
	category string,
) (int, error) {

	all := []model.Event{
		{ID: "1", Title: "Rock Concert", Category: "concert"},
		{ID: "2", Title: "Cinema Night", Category: "film"},
		{ID: "3", Title: "Business Seminar", Category: "seminar"},
	}

	count := 0
	for _, e := range all {
		if category == "" || strings.EqualFold(e.Category, category) {
			count++
		}
	}

	return count, nil
}

func (r *InMemoryEventRepository) Create(
	ctx context.Context,
	event model.Event,
) (model.Event, error) {
	// in-memory: just echo back
	return event, nil
}
