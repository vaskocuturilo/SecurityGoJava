package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang/internal/config"
	"golang/model"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) SignUp(ctx context.Context, credential *model.Credential) error {

	hashedPassword, err := config.HashPassword(credential.Password)

	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
        INSERT INTO users_db (id, username, email, password, enabled, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err = r.db.ExecContext(ctx, query,
		uuid.New(),
		credential.UserName,
		credential.UserName,
		string(hashedPassword),
		true,
		time.Now().UTC().UnixMilli(),
		time.Now().UTC().UnixMilli(),
	)
	return err
}

func (r *PostgresUserRepository) Login(ctx context.Context, username, password string) (*model.User, error) {
	var user model.User
	var hashedPwd string

	query := `SELECT id, username, password FROM users_db WHERE username = $1`

	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &hashedPwd)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))

	if err != nil {
		return nil, model.ErrInvalidCredentials
	}

	return &user, nil
}

func (r *PostgresUserRepository) Refresh(ctx context.Context, request *model.RefreshRequest) error {
	return errors.New("not implemented")
}
