package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserStore {
	return &store{pg}
}

func (s *store) GetUserByUsername(ctx context.Context, username string) (user *core.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user = new(core.User)

	stmt := `SELECT id, username, password_hash FROM users WHERE username = $1`
	err = s.DB.QueryRowContext(ctx, stmt, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrInvalidCredentials
		}
		return nil, err
	}

	return user, nil
}

func (s *store) AddUser(ctx context.Context, user core.User) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// starting transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stmt := `SELECT id FROM users
	WHERE username = $1`

	var userID int
	err = tx.QueryRowContext(ctx, stmt, user.Username).Scan(&userID)
	if userID != 0 {
		return 0, core.ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	stmt = `INSERT INTO users (username, password_hash)
	VALUES ($1, $2) RETURNING id`

	err = tx.QueryRowContext(ctx, stmt, user.Username, user.PasswordHash).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *store) UpdateUser(ctx context.Context, user core.User) (userID int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `UPDATE users SET password_hash = $1 WHERE id = $2 RETURNING id`
	err = s.DB.QueryRowContext(ctx, stmt, user.PasswordHash, user.ID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
