package identity

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateInitialOwner(ctx context.Context, email, passwordHash, storeName string) (*User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(702184031)`); err != nil {
		return nil, err
	}

	var hasUsers bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users)`).Scan(&hasUsers); err != nil {
		return nil, err
	}
	if hasUsers {
		return nil, ErrInitialOwnerExists
	}

	var storeID uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT id
		FROM stores
		WHERE is_legacy = TRUE AND owner_user_id IS NULL
		ORDER BY created_at
		LIMIT 1
		FOR UPDATE
	`).Scan(&storeID)
	if errors.Is(err, pgx.ErrNoRows) {
		storeID = uuid.New()
		_, err = tx.Exec(ctx, `INSERT INTO stores (id, name) VALUES ($1, $2)`, storeID, storeName)
	} else if err == nil {
		_, err = tx.Exec(ctx, `UPDATE stores SET name = $1 WHERE id = $2`, storeName, storeID)
	}
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:        uuid.New(),
		StoreID:   storeID,
		Email:     email,
		Role:      RoleOwner,
		CreatedAt: time.Now().UTC(),
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, store_id, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, user.ID, user.StoreID, user.Email, passwordHash, user.Role, user.CreatedAt)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, `UPDATE stores SET owner_user_id = $1 WHERE id = $2`, user.ID, storeID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) FindUserByEmail(ctx context.Context, email string) (*User, string, error) {
	var user User
	var passwordHash string
	err := r.db.QueryRow(ctx, `
		SELECT id, store_id, email, role, password_hash, created_at
		FROM users
		WHERE LOWER(email) = $1 AND disabled_at IS NULL
	`, email).Scan(&user.ID, &user.StoreID, &user.Email, &user.Role, &passwordHash, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, "", ErrInvalidCredentials
	}
	if err != nil {
		return nil, "", err
	}
	return &user, passwordHash, nil
}

func (r *PostgresRepository) CreateSession(ctx context.Context, tokenHash, csrfHash []byte, userID uuid.UUID, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sessions (token_hash, user_id, csrf_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`, tokenHash, userID, csrfHash, expiresAt)
	return err
}

func (r *PostgresRepository) GetSession(ctx context.Context, tokenHash []byte) (*Session, error) {
	var session Session
	err := r.db.QueryRow(ctx, `
		SELECT s.user_id, u.store_id, u.role, s.csrf_hash, s.expires_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1
		  AND s.expires_at > NOW()
		  AND u.disabled_at IS NULL
	`, tokenHash).Scan(&session.UserID, &session.StoreID, &session.Role, &session.CSRFHash, &session.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrInvalidSession
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *PostgresRepository) DeleteSession(ctx context.Context, tokenHash []byte) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, tokenHash)
	return err
}
