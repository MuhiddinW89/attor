CREATE TABLE stores (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_legacy BOOLEAN NOT NULL DEFAULT FALSE,
    owner_user_id UUID UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO stores (id, name, is_legacy)
VALUES ('00000000-0000-0000-0000-000000000001', 'Migrated store', TRUE);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    email VARCHAR(320) NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(16) NOT NULL CHECK (role IN ('owner', 'employee')),
    disabled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (id, store_id)
);

CREATE UNIQUE INDEX users_email_lower_idx ON users (LOWER(email));
CREATE UNIQUE INDEX users_one_owner_per_store_idx
    ON users (store_id) WHERE role = 'owner';

ALTER TABLE stores
    ADD CONSTRAINT fk_stores_owner_same_store
    FOREIGN KEY (owner_user_id, id) REFERENCES users(id, store_id);

CREATE TABLE sessions (
    token_hash BYTEA PRIMARY KEY CHECK (OCTET_LENGTH(token_hash) = 32),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    csrf_hash BYTEA NOT NULL CHECK (OCTET_LENGTH(csrf_hash) = 32),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (expires_at > created_at)
);

CREATE INDEX sessions_user_id_idx ON sessions(user_id);
CREATE INDEX sessions_expires_at_idx ON sessions(expires_at);

CREATE TABLE employee_invites (
    token_hash BYTEA PRIMARY KEY CHECK (OCTET_LENGTH(token_hash) = 32),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    email VARCHAR(320) NOT NULL,
    invited_by_user_id UUID NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    accepted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (invited_by_user_id, store_id)
        REFERENCES users(id, store_id) ON DELETE CASCADE,
    CHECK (expires_at > created_at)
);

CREATE INDEX employee_invites_store_id_idx ON employee_invites(store_id);
CREATE INDEX employee_invites_expires_at_idx ON employee_invites(expires_at);

ALTER TABLE clients
    ADD COLUMN store_id UUID NOT NULL
    DEFAULT '00000000-0000-0000-0000-000000000001';
ALTER TABLE sales
    ADD COLUMN store_id UUID NOT NULL
    DEFAULT '00000000-0000-0000-0000-000000000001';
ALTER TABLE reminders
    ADD COLUMN store_id UUID NOT NULL
    DEFAULT '00000000-0000-0000-0000-000000000001';

ALTER TABLE clients
    ADD CONSTRAINT fk_clients_store FOREIGN KEY (store_id)
    REFERENCES stores(id) ON DELETE CASCADE,
    ADD CONSTRAINT clients_id_store_id_unique UNIQUE (id, store_id);

ALTER TABLE sales
    ADD CONSTRAINT fk_sales_store FOREIGN KEY (store_id)
    REFERENCES stores(id) ON DELETE CASCADE,
    ADD CONSTRAINT sales_id_store_id_unique UNIQUE (id, store_id),
    ADD CONSTRAINT fk_sales_client_store FOREIGN KEY (client_id, store_id)
    REFERENCES clients(id, store_id) ON DELETE CASCADE;

ALTER TABLE reminders
    ADD CONSTRAINT fk_reminders_store FOREIGN KEY (store_id)
    REFERENCES stores(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_reminders_client_store FOREIGN KEY (client_id, store_id)
    REFERENCES clients(id, store_id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_reminders_sale_store FOREIGN KEY (sale_id, store_id)
    REFERENCES sales(id, store_id) ON DELETE CASCADE;

CREATE INDEX clients_store_id_idx ON clients(store_id);
CREATE INDEX sales_store_id_idx ON sales(store_id);
CREATE INDEX reminders_store_id_idx ON reminders(store_id);