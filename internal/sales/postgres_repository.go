package sales

import (
	"context"
	"database/sql"
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
	sale *Sale,
) error {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return identity.ErrInvalidSession
	}
	sale.StoreID = session.StoreID
	sale.CreatedByUserID = session.UserID

	const query = `
		INSERT INTO sales (
			id,
			store_id,
			created_by_user_id,
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
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
		)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		sale.ID,
		sale.StoreID,
		sale.CreatedByUserID,
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
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return nil, identity.ErrInvalidSession
	}

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
		WHERE id = $1 AND store_id = $2
		  AND ($3 = 'owner' OR created_by_user_id = $4)
	`

	var sale Sale
	var comment sql.NullString

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		session.StoreID,
		session.Role,
		session.UserID,
	).Scan(
		&sale.ID,
		&sale.ClientID,
		&sale.PerfumeName,
		&sale.VolumeML,
		&sale.Price,
		&comment,
		&sale.SaleDate,
		&sale.CreatedAt,
		&sale.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSaleNotFound
		}
		return nil, err
	}

	if comment.Valid {
		sale.Comment = &comment.String
	}

	return &sale, nil
}

func (r *PostgresRepository) ListByClientID(
	ctx context.Context,
	clientID uuid.UUID,
) ([]*Sale, error) {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return nil, identity.ErrInvalidSession
	}

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
		WHERE client_id = $1 AND store_id = $2
		  AND ($3 = 'owner' OR created_by_user_id = $4)
		ORDER BY sale_date DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		clientID,
		session.StoreID,
		session.Role,
		session.UserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []*Sale

	for rows.Next() {

		var sale Sale
		var comment sql.NullString

		err := rows.Scan(
			&sale.ID,
			&sale.ClientID,
			&sale.PerfumeName,
			&sale.VolumeML,
			&sale.Price,
			&comment,
			&sale.SaleDate,
			&sale.CreatedAt,
			&sale.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if comment.Valid {
			sale.Comment = &comment.String
		}

		sales = append(sales, &sale)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sales, nil
}
