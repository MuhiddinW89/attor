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

func (r *PostgresRepository) Create(ctx context.Context, client *Client) error {

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

func (r *PostgresRepository) GetByPhone(ctx context.Context, phone string) (*Client, error) {

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

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Client, error) {
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

func (r *PostgresRepository) List(ctx context.Context) ([]*Client, error) {

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
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
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
		if err != nil {
			return nil, err
		}

		clients = append(
			clients,
			&client,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}
