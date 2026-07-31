CREATE TABLE clients (
    id UUID PRIMARY KEY,

    full_name VARCHAR(255) NOT NULL,

    phone VARCHAR(20) NOT NULL UNIQUE,

    instagram VARCHAR(255),

    birth_date DATE,

    created_at TIMESTAMP NOT NULL,

    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_clients_phone
ON clients(phone);