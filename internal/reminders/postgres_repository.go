package reminders

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
			is_sent,
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
		reminder.IsSent,
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
			is_sent,
			created_at,
			updated_at
		FROM reminders
		WHERE is_sent = false
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
			&reminder.IsSent,
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
			is_sent = true,
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

func (r *PostgresRepository) GetNearestByClientID(
	ctx context.Context,
	clientID uuid.UUID,
) (*Reminder, error) {

	query := `
		SELECT
			id,
			client_id,
			sale_id,
			reminder_at,
			is_sent,
			created_at,
			updated_at
		FROM reminders
		WHERE client_id = $1
		ORDER BY reminder_at
		LIMIT 1
	`

	var reminder Reminder

	err := r.db.QueryRow(ctx, query, clientID).Scan(
		&reminder.ID,
		&reminder.ClientID,
		&reminder.SaleID,
		&reminder.ReminderAt,
		&reminder.IsSent,
		&reminder.CreatedAt,
		&reminder.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrReminderNotFound
		}

		return nil, err
	}

	return &reminder, nil
}
