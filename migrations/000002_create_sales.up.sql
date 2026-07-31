CREATE TABLE sales (
    id UUID PRIMARY KEY,

    client_id UUID NOT NULL,

    perfume_name VARCHAR(255) NOT NULL,

    volume_ml INTEGER NOT NULL,

    price NUMERIC(12,2) NOT NULL,

    comment TEXT,

    sale_date TIMESTAMP NOT NULL,

    created_at TIMESTAMP NOT NULL,

    updated_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_sales_client
        FOREIGN KEY (client_id)
        REFERENCES clients(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_sales_client_id
ON sales(client_id);

CREATE INDEX idx_sales_sale_date
ON sales(sale_date);