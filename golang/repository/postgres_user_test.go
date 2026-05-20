package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresUserRepository_SignUp_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()

	defer db.Close()

	repo := NewPostgresUserRepository(db)

	user := &model.Credential{Email: "Title@title.com", Password: "Description"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users_db (id, email, password, enabled, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)")).
		WithArgs(sqlmock.AnyArg(), user.Email, user.Password, true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.SignUp(context.Background(), user)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresUserRepository_GetByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	expectedId := "123"
	expectedEmail := "Email@email.com"
	expectedPassword := "hashed-password"

	rows := sqlmock.NewRows([]string{"expectedId", "email", "password"}).
		AddRow(expectedId, expectedEmail, expectedPassword)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password FROM users_db WHERE email = $1")).
		WithArgs(expectedEmail).
		WillReturnRows(rows)

	user, err := repo.GetByEmail(context.Background(), expectedEmail)

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	if user == nil {
		t.Fatal("expected user to be not nil")
	}

	if user.ID != expectedId {
		t.Errorf("expected expectedId '123', got %s", user.ID)
	}

	if user.Email != expectedEmail {
		t.Errorf("expected email 'Email@email.com', got %s", user.Email)
	}

	if user.HashedPassword != expectedPassword {
		t.Errorf("expected password 'hashed-password', got %s", user.HashedPassword)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresUserRepository_GetByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresUserRepository(db)

	expectedEmail := "unknow@email.com"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password FROM users_db WHERE email = $1")).
		WithArgs(expectedEmail).WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByEmail(context.Background(), expectedEmail)

	if !errors.Is(err, model.ErrUserNotFound) {
		t.Errorf("expected error %v, got %v", model.ErrUserNotFound, err)
	}

	if user != nil {
		t.Errorf("expected user to be nil, got %v", user)
	}
}
