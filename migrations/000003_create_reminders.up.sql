CREATE TABLE reminders (
    id UUID PRIMARY KEY,

    client_id UUID NOT NULL,

    sale_id UUID NOT NULL,

    reminder_at TIMESTAMP NOT NULL,

    is_sent BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL,

    updated_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_reminders_client
        FOREIGN KEY (client_id)
        REFERENCES clients(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_reminders_sale
        FOREIGN KEY (sale_id)
        REFERENCES sales(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_reminders_client_id
ON reminders(client_id);

CREATE INDEX idx_reminders_sale_id
ON reminders(sale_id);

CREATE INDEX idx_reminders_reminder_at
ON reminders(reminder_at);

CREATE INDEX idx_reminders_is_sent
ON reminders(is_sent);