package model

type TicketStatus string

const (
	StatusAvailable TicketStatus = "AVAILABLE"
	StatusReserved  TicketStatus = "RESERVED"
	StatusSold      TicketStatus = "SOLD"
)
