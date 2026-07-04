package sales

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(
	db *pgxpool.Pool,
) *PostgresRepository {
	return &PostgresRepository {
		db:db,
	}
}


func (r *PostgresRepository) Create (
	ctx context.Context,
	sale *Sale,
)error {
	const query = `
		INSERT INTO sales (
		id,
		client_id,
		perfume_name,
		volume_ml,
		price,
		comment,
		sale_date,
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
		sale.ID,
		sale.ClientID,
		sale.PerfumeName,
		sale.VolumeML,
		sale.Price,
		sale.Comment,
		sale.SaleDate,
		sale.CreatedAt,
		sale.UpdatedAt,
	)

	return err
}

func (r *PostgresRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Sale, error) {
	const query = `
		SELECT 
			id,
			client_id,
			perfume_name,
			volume_ml,
			price,
			comment,
			sale_date,
			created_at,
			updated_at
		FROM sales
		WHERE id = $1	
	`

	var sale Sale

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&sale.ClientID
	)

	return &sale, nil
}