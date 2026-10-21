DROP INDEX reminders_creator_store_idx;
DROP INDEX sales_creator_store_date_idx;
DROP INDEX clients_creator_store_idx;

ALTER TABLE reminders DROP CONSTRAINT fk_reminders_creator_store, DROP COLUMN created_by_user_id;
ALTER TABLE sales DROP CONSTRAINT fk_sales_creator_store, DROP COLUMN created_by_user_id;
ALTER TABLE clients
    DROP CONSTRAINT fk_clients_creator_store,
    DROP CONSTRAINT clients_store_phone_unique,
    DROP COLUMN created_by_user_id,
    ADD CONSTRAINT clients_phone_key UNIQUE (phone);