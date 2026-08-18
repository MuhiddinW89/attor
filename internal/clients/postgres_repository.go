package clients

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(
	db *pgxpool.Pool,
) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(
	ctx context.Context,
	client *Client,
) error {

	query := `
		INSERT INTO clients (
			id,
			full_name,
			phone,
			instagram,
			birth_date,
			created_at,
			updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7
		)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		client.ID,
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
		WHERE phone = $1
	`

	var client Client

	err := r.db.QueryRow(
		ctx,
		query,
		phone,
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
		WHERE id = $1
	`

	var client Client

	err := r.db.QueryRow(
		ctx,
		query,
		id,
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
		WHERE
			$1 = ''
			OR full_name ILIKE '%' || $1 || '%'
			OR phone ILIKE '%' || $1 || '%'
		ORDER BY full_name ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
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
			ON s.client_id = c.id
		WHERE
			$1 = ''
			OR c.full_name ILIKE '%' || $1 || '%'
			OR c.phone ILIKE '%' || $1 || '%'
			OR COALESCE(c.instagram, '') ILIKE '%' || $1 || '%'
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
