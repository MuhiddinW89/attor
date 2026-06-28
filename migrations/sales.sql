CREATE TABLE sales (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL REFERENCES clients(id),
    perfume_name VARCHAR(255) NOT NULL,
    volume_ml INTEGER NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    comment TEXT,
    sale_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_sales_client_id
ON sales(client_id);