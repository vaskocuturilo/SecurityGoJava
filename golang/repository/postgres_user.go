package repository

import (
	"context"
	"database/sql"
	"golang/model"
	"time"

	"github.com/google/uuid"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) SignUp(ctx context.Context, credential *model.Credential) error {
	query := `
        INSERT INTO users_db (id, username, email, password, enabled, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.ExecContext(ctx, query,
		uuid.New(), credential.UserName, credential.UserName, credential.Password, true, time.Now().UTC().UnixMilli(), time.Now().UTC().UnixMilli())
	return err
}

func (r *PostgresUserRepository) Login(ctx context.Context, credential *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresUserRepository) Refresh(ctx context.Context, request *model.RefreshRequest) error {
	//TODO implement me
	panic("implement me")
}
