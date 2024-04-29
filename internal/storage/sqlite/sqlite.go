package sqlite

import (
	"SSO/internal/domain/models"
	"SSO/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

const (
	emptyValue = 0
)

var (
	emptyUser = &models.User{}
)

type Storage struct {
	db *sql.DB
}

// Connect connects to db, returns Storage struct
func Connect(storagePath string) (*Storage, error) {
	const op = "sqlite.Connect"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveUser method of Storage.
// Implements saving user into database sqlite
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "sqlite.SaveUser"

	// Prepare instead of exec because of better performance first one
	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return emptyValue, fmt.Errorf("%s : %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteError sqlite3.Error

		if errors.As(err, &sqliteError) && errors.Is(sqliteError, sqlite3.ErrConstraintUnique) {
			return emptyValue, fmt.Errorf("%s, %w", op, storage.ErrUserAlreadyExists)
		}

		return emptyValue, fmt.Errorf("%s : %w", op, err)
	}

	newID, err := res.LastInsertId()
	if err != nil {
		return emptyValue, fmt.Errorf("%s : %w", op, err)
	}

	return newID, nil
}

// User method of Storage.
// Implements getting user from database
func (s *Storage) User(ctx context.Context, email string) (user *models.User, err error) {
	const op = "sqlite.User"

	// prepare for exec
	stmt, err := s.db.Prepare("SELECT * FROM users WHERE (email) VALUES (?)")
	if err != nil {
		return emptyUser, fmt.Errorf("%s : %w", op, err)
	}

	// exec statement context
	row := stmt.QueryRowContext(ctx, email)

	user = &models.User{}
	if err := row.Scan(&user.Email, &user.ID, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return emptyUser, fmt.Errorf("%s : %w", op, storage.ErrUserNotFound)
		}

		return emptyUser, fmt.Errorf("%s : %w", op, err)
	}

	return user, nil
}
