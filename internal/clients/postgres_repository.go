package clients

import (
	"context"
	"errors"

	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/MuhiddinW89/attor/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	db database.Querier
}

func NewPostgresRepository(
	db database.Querier,
) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(
	ctx context.Context,
	client *Client,
) error {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return identity.ErrInvalidSession
	}
	client.StoreID = session.StoreID
	client.CreatedByUserID = &session.UserID

	query := `
		INSERT INTO clients (
			id,
			store_id,
			created_by_user_id,
			full_name,
			phone,
			instagram,
			birth_date,
			created_at,
			updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9
		)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		client.ID,
		client.StoreID,
		client.CreatedByUserID,
		client.FullName,
		client.Phone,
		client.Instagram,
		client.BirthDate,
		client.CreatedAt,
		client.UpdatedAt,
	)

	return err
}

func (r *PostgresRepository) GetByPhone(
	ctx context.Context,
	phone string,
) (*Client, error) {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return nil, identity.ErrInvalidSession
	}

	query := `
		SELECT
			id,
			full_name,
			phone,
			instagram,
			birth_date,
			created_at,
			updated_at
		FROM clients
		WHERE store_id = $1
		  AND phone = $2
		  AND ($3 = 'owner' OR created_by_user_id = $4 OR EXISTS (
			SELECT 1 FROM sales s
			WHERE s.client_id = clients.id
			  AND s.store_id = clients.store_id
			  AND s.created_by_user_id = $4
		  ))
	`

	var client Client

	err := r.db.QueryRow(
		ctx,
		query,
		session.StoreID,
		phone,
		session.Role,
		session.UserID,
	).Scan(
		&client.ID,
		&client.FullName,
		&client.Phone,
		&client.Instagram,
		&client.BirthDate,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClientNotFound
		}

		return nil, err
	}

	return &client, nil
}

func (r *PostgresRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Client, error) {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return nil, identity.ErrInvalidSession
	}

	query := `
		SELECT
			id,
			full_name,
			phone,
			instagram,
			birth_date,
			created_at,
			updated_at
		FROM clients
		WHERE id = $1 AND store_id = $2
		  AND ($3 = 'owner' OR created_by_user_id = $4 OR EXISTS (
			SELECT 1 FROM sales s
			WHERE s.client_id = clients.id
			  AND s.store_id = clients.store_id
			  AND s.created_by_user_id = $4
		  ))
	`

	var client Client

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		session.StoreID,
		session.Role,
		session.UserID,
	).Scan(
		&client.ID,
		&client.FullName,
		&client.Phone,
		&client.Instagram,
		&client.BirthDate,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClientNotFound
		}

		return nil, err
	}

	return &client, nil
}

func (r *PostgresRepository) List(
	ctx context.Context,
	search string,
) ([]*Client, error) {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return nil, identity.ErrInvalidSession
	}

	const query = `
		SELECT
			id,
			full_name,
			phone,
			instagram,
			birth_date,
			created_at,
			updated_at
		FROM clients
		WHERE store_id = $1
		  AND ($2 = 'owner' OR created_by_user_id = $3 OR EXISTS (
			SELECT 1 FROM sales s
			WHERE s.client_id = clients.id
			  AND s.store_id = clients.store_id
			  AND s.created_by_user_id = $3
		  ))
		  AND ($4 = '' OR full_name ILIKE '%' || $4 || '%' OR phone ILIKE '%' || $4 || '%')
		ORDER BY full_name ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		session.StoreID,
		session.Role,
		session.UserID,
		search,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*Client

	for rows.Next() {

		var client Client

		err := rows.Scan(
			&client.ID,
			&client.FullName,
			&client.Phone,
			&client.Instagram,
			&client.BirthDate,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
		//fmt.Printf("DEBUG FULL NAME: %q\n", client.FullName)
		if err != nil {
			return nil, err
		}

		clients = append(clients, &client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *PostgresRepository) ListWithStats(
	ctx context.Context,
	search string,
) ([]*ClientListStatsItem, error) {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return nil, identity.ErrInvalidSession
	}

	const query = `
		SELECT
			c.id,
			c.full_name,
			c.phone,
			c.instagram,
			COUNT(s.id) AS total_sales,
			COALESCE(SUM(s.price), 0) AS total_amount
		FROM clients c
		LEFT JOIN sales s
			ON s.client_id = c.id AND s.store_id = c.store_id
			AND ($2 = 'owner' OR s.created_by_user_id = $3)
		WHERE
			c.store_id = $1
			AND ($2 = 'owner' OR c.created_by_user_id = $3 OR EXISTS (
				SELECT 1 FROM sales own_sales
				WHERE own_sales.client_id = c.id AND own_sales.store_id = c.store_id
				  AND own_sales.created_by_user_id = $3
			))
			AND ($4 = '' OR c.full_name ILIKE '%' || $4 || '%'
				OR c.phone ILIKE '%' || $4 || '%'
				OR COALESCE(c.instagram, '') ILIKE '%' || $4 || '%')
		GROUP BY
			c.id,
			c.full_name,
			c.phone,
			c.instagram
		ORDER BY
			c.full_name ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		session.StoreID,
		session.Role,
		session.UserID,
		search,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*ClientListStatsItem

	for rows.Next() {

		var client ClientListStatsItem

		err := rows.Scan(
			&client.ID,
			&client.FullName,
			&client.Phone,
			&client.Instagram,
			&client.TotalSales,
			&client.TotalAmount,
		)
		if err != nil {
			return nil, err
		}

		clients = append(clients, &client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}
