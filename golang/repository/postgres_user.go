package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang/model"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) SignUp(ctx context.Context, credential *model.Credential) error {

	query := `
        INSERT INTO users_db (id, email, password, enabled, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	uid := uuid.New().String()
	email := credential.Email
	pass := credential.Password
	now := time.Now().UTC()

	_, err := r.db.ExecContext(ctx, query,
		uid,
		email,
		pass,
		true,
		now,
		now,
	)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrAlreadyExists
		}
	}
	return err
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, password, role FROM users_db WHERE email = $1`

	var u model.User
	var hashedPwd string

	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &hashedPwd, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	u.HashedPassword = hashedPwd
	return &u, nil
}

func (r *PostgresUserRepository) Refresh(refreshToken string) (string, string, error) {
	return "", "", errors.New("not implemented")
}
