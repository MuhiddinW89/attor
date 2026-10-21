package reminders

import (
	"context"
	"errors"
	"fmt"

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
	reminder *Reminder,
) error {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return identity.ErrInvalidSession
	}
	reminder.StoreID = session.StoreID
	reminder.CreatedByUserID = session.UserID
	const query = `
		INSERT INTO reminders (
			id,
			store_id,
			created_by_user_id,
			client_id,
			sale_id,
			reminder_at,
			is_sent,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		reminder.ID,
		reminder.StoreID,
		reminder.CreatedByUserID,
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
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return nil, identity.ErrInvalidSession
	}
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
		WHERE is_sent = false AND store_id = $1
		  AND ($2 = 'owner' OR created_by_user_id = $3)
		ORDER BY reminder_at ASC
	`

	rows, err := r.db.Query(ctx, query, session.StoreID, session.Role, session.UserID)
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

		reminders = append(
			reminders,
			&reminder,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *PostgresRepository) GetNearestByClientID(
	ctx context.Context,
	clientID uuid.UUID,
) (*Reminder, error) {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return nil, identity.ErrInvalidSession
	}
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
		WHERE client_id = $1 AND store_id = $2
		  AND ($3 = 'owner' OR created_by_user_id = $4)
		ORDER BY reminder_at
		LIMIT 1
	`

	var reminder Reminder

	err := r.db.QueryRow(
		ctx,
		query,
		clientID,
		session.StoreID,
		session.Role,
		session.UserID,
	).Scan(
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

func (r *PostgresRepository) MarkSent(
	ctx context.Context,
	id uuid.UUID,
) error {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return identity.ErrInvalidSession
	}
	const query = `
		UPDATE reminders
		SET
			is_sent = true,
			updated_at = NOW()
		WHERE id = $1 AND store_id = $2
		  AND ($3 = 'owner' OR created_by_user_id = $4)
	`

	commandTag, err := r.db.Exec(
		ctx,
		query,
		id,
		session.StoreID,
		session.Role,
		session.UserID,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrReminderNotFound
	}

	return nil
}

// ListDetailed returns reminders filtered by status.
//
// Supported statuses:
//   - all
//   - urgent
//   - overdue
//   - completed
func (r *PostgresRepository) ListDetailed(
	ctx context.Context,
	status string,
) ([]*ReminderListItem, error) {
	var query string

	switch status {
	case "all":
		query = `
			SELECT
				r.id,
				r.client_id,
				c.full_name,
				c.phone,
				s.perfume_name,
				s.volume_ml,
				s.price,
				COALESCE(s.comment, '') AS comment,
				s.sale_date,
				r.reminder_at,
				r.is_sent
			FROM reminders r
			INNER JOIN clients c
				ON c.id = r.client_id
			INNER JOIN sales s
				ON s.id = r.sale_id
			WHERE r.store_id = $1 AND ($2 = 'owner' OR r.created_by_user_id = $3)
			ORDER BY
				r.is_sent ASC,
				r.reminder_at ASC
		`

	case "urgent":
		query = `
			SELECT
				r.id,
				r.client_id,
				c.full_name,
				c.phone,
				s.perfume_name,
				s.volume_ml,
				s.price,
				COALESCE(s.comment, '') AS comment,
				s.sale_date,
				r.reminder_at,
				r.is_sent
			FROM reminders r
			INNER JOIN clients c
				ON c.id = r.client_id
			INNER JOIN sales s
				ON s.id = r.sale_id
			WHERE
				r.is_sent = false
				AND r.store_id = $1
				AND ($2 = 'owner' OR r.created_by_user_id = $3)
				AND r.reminder_at >= NOW()
				AND r.reminder_at <= NOW() + INTERVAL '3 days'
			ORDER BY r.reminder_at ASC
		`

	case "overdue":
		query = `
			SELECT
				r.id,
				r.client_id,
				c.full_name,
				c.phone,
				s.perfume_name,
				s.volume_ml,
				s.price,
				COALESCE(s.comment, '') AS comment,
				s.sale_date,
				r.reminder_at,
				r.is_sent
			FROM reminders r
			INNER JOIN clients c
				ON c.id = r.client_id
			INNER JOIN sales s
				ON s.id = r.sale_id
			WHERE
				r.is_sent = false
				AND r.store_id = $1
				AND ($2 = 'owner' OR r.created_by_user_id = $3)
				AND r.reminder_at < NOW()
			ORDER BY r.reminder_at ASC
		`

	case "completed":
		query = `
			SELECT
				r.id,
				r.client_id,
				c.full_name,
				c.phone,
				s.perfume_name,
				s.volume_ml,
				s.price,
				COALESCE(s.comment, '') AS comment,
				s.sale_date,
				r.reminder_at,
				r.is_sent
			FROM reminders r
			INNER JOIN clients c
				ON c.id = r.client_id
			INNER JOIN sales s
				ON s.id = r.sale_id
			WHERE r.is_sent = true
			  AND r.store_id = $1
			  AND ($2 = 'owner' OR r.created_by_user_id = $3)
			ORDER BY r.updated_at DESC
		`

	default:
		return nil, fmt.Errorf(
			"invalid reminder status: %s",
			status,
		)
	}

	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return nil, identity.ErrInvalidSession
	}
	rows, err := r.db.Query(ctx, query, session.StoreID, session.Role, session.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []*ReminderListItem

	for rows.Next() {
		var item ReminderListItem

		err := rows.Scan(
			&item.ID,
			&item.ClientID,
			&item.ClientName,
			&item.Phone,
			&item.PerfumeName,
			&item.VolumeML,
			&item.Price,
			&item.Comment,
			&item.SaleDate,
			&item.ReminderAt,
			&item.IsSent,
		)
		if err != nil {
			return nil, err
		}

		reminders = append(
			reminders,
			&item,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reminders, nil
}

// ListPendingDetailed is kept for compatibility
// with the existing service.
func (r *PostgresRepository) ListPendingDetailed(
	ctx context.Context,
) ([]*ReminderListItem, error) {
	return r.ListDetailed(ctx, "all")
}
