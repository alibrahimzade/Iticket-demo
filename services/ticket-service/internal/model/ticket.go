package model

import "time"

type Ticket struct {
	ID        string
	EventID   string
	Zone      string
	Price     float64
	Curreny   string
	Status    string
	CreatedAt time.Time
}
