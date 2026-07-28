package reminders

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(
	ctx context.Context,
	reminder *Reminder,
) error {

	const query = `
		INSERT INTO reminders (
			id,
			client_id,
			sale_id,
			reminder_at,
			is_completed,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		reminder.ID,
		reminder.ClientID,
		reminder.SaleID,
		reminder.ReminderAt,
		reminder.IsCompleted,
		reminder.CreatedAt,
		reminder.UpdatedAt,
	)

	return err
}

func (r *PostgresRepository) ListPending(
	ctx context.Context,
) ([]*Reminder, error) {

	const query = `
		SELECT
			id,
			client_id,
			sale_id,
			reminder_at,
			is_completed,
			created_at,
			updated_at
		FROM reminders
		WHERE is_completed = false
		ORDER BY reminder_at ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []*Reminder

	for rows.Next() {
		var reminder Reminder

		err := rows.Scan(
			&reminder.ID,
			&reminder.ClientID,
			&reminder.SaleID,
			&reminder.ReminderAt,
			&reminder.IsCompleted,
			&reminder.CreatedAt,
			&reminder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		reminders = append(reminders, &reminder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *PostgresRepository) MarkCompleted(
	ctx context.Context,
	id uuid.UUID,
) error {

	const query = `
		UPDATE reminders
		SET
			is_completed = true,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	return err
}
