CREATE TABLE IF NOT EXISTS tickets (
                                       id UUID PRIMARY KEY,
                                       event_id UUID NOT NULL,
                                       zone TEXT NOT NULL,
                                       price NUMERIC(10,2) NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_tickets_event_id
    ON tickets(event_id);
