ALTER TABLE clients
    DROP CONSTRAINT clients_phone_key,
    ADD CONSTRAINT clients_store_phone_unique UNIQUE (store_id, phone),
    ADD COLUMN created_by_user_id UUID;

ALTER TABLE sales
    ADD COLUMN created_by_user_id UUID;

ALTER TABLE reminders
    ADD COLUMN created_by_user_id UUID;

ALTER TABLE clients
    ADD CONSTRAINT fk_clients_creator_store
    FOREIGN KEY (created_by_user_id, store_id) REFERENCES users(id, store_id);
ALTER TABLE sales
    ADD CONSTRAINT fk_sales_creator_store
    FOREIGN KEY (created_by_user_id, store_id) REFERENCES users(id, store_id);
ALTER TABLE reminders
    ADD CONSTRAINT fk_reminders_creator_store
    FOREIGN KEY (created_by_user_id, store_id) REFERENCES users(id, store_id);

CREATE INDEX clients_creator_store_idx ON clients(created_by_user_id, store_id);
CREATE INDEX sales_creator_store_date_idx ON sales(created_by_user_id, store_id, sale_date);
CREATE INDEX reminders_creator_store_idx ON reminders(created_by_user_id, store_id);