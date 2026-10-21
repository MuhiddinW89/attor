DROP TABLE employee_invites;
DROP TABLE sessions;

ALTER TABLE reminders
    DROP CONSTRAINT fk_reminders_sale_store,
    DROP CONSTRAINT fk_reminders_client_store,
    DROP CONSTRAINT fk_reminders_store;
ALTER TABLE sales
    DROP CONSTRAINT fk_sales_client_store,
    DROP CONSTRAINT sales_id_store_id_unique,
    DROP CONSTRAINT fk_sales_store;
ALTER TABLE clients
    DROP CONSTRAINT clients_id_store_id_unique,
    DROP CONSTRAINT fk_clients_store;

DROP INDEX reminders_store_id_idx;
DROP INDEX sales_store_id_idx;
DROP INDEX clients_store_id_idx;

ALTER TABLE reminders DROP COLUMN store_id;
ALTER TABLE sales DROP COLUMN store_id;
ALTER TABLE clients DROP COLUMN store_id;

ALTER TABLE stores DROP CONSTRAINT fk_stores_owner_same_store;
DROP TABLE users;
DROP TABLE stores;