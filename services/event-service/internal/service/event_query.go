package service

type EventQuery struct {
	Category string
	Page     int
	Size     int
	SortBy   string // title | category | id
	Order    string // asc | desc
}
