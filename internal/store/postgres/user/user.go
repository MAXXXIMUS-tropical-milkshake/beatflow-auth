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

func (s *store) GetUserByEmail(ctx context.Context, email string) (user *core.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user = new(core.User)

	stmt := `SELECT id, username, email, password_hash FROM users WHERE email = $1`
	err = s.DB.QueryRowContext(ctx, stmt, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrInvalidCredentials
		}
		return nil, err
	}

	return user, nil
}

func (s *store) GetUserByUsername(ctx context.Context, username string) (user *core.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user = new(core.User)

	stmt := `SELECT id, username, email, password_hash FROM users WHERE username = $1`
	err = s.DB.QueryRowContext(ctx, stmt, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrInvalidCredentials
		}
		return nil, err
	}

	return user, nil
}

func (s *store) GetUserByID(ctx context.Context, userID int) (user *core.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user = new(core.User)

	stmt := `SELECT id, username, email, password_hash FROM users WHERE id = $1`
	err = s.DB.QueryRowContext(ctx, stmt, userID).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *store) AddUser(ctx context.Context, user core.User) (userID int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// starting transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback() // nolint
		} else {
			tx.Commit() // nolint
		}
	}()

	stmt := `SELECT id FROM users
	WHERE username = $1`

	err = tx.QueryRowContext(ctx, stmt, user.Username).Scan(&userID)
	if userID != 0 {
		return 0, core.ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	stmt = `INSERT INTO users (username, email, password_hash)
	VALUES ($1, $2, $3) RETURNING id`

	err = tx.QueryRowContext(ctx, stmt, user.Username, user.Email, user.PasswordHash).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, err
}

func (s *store) UpdateUser(ctx context.Context, user core.UpdateUser) (userID int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var password *string
	if user.Password != nil {
		password = &user.Password.NewPassword
	}

	// starting transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback() // nolint
		} else {
			tx.Commit() // nolint
		}
	}()

	if user.Email != nil {
		stmt := `SELECT id FROM users WHERE email = $1`
		err = tx.QueryRowContext(ctx, stmt, user.Email).Scan(&userID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		} else if userID != 0 {
			return 0, core.ErrEmailAlreadyExists
		}
	}

	if user.Username != nil {
		stmt := `SELECT id FROM users WHERE username = $1`
		err = tx.QueryRowContext(ctx, stmt, user.Username).Scan(&userID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		} else if userID != 0 {
			return 0, core.ErrUsernameAlreadyExists
		}
	}

	stmt := `UPDATE users SET
	password_hash = COALESCE($1, password_hash),
	username = COALESCE($2, username),
	email = COALESCE($3, email)
	WHERE id = $4
	RETURNING id`
	err = tx.QueryRowContext(ctx, stmt, password, user.Username, user.Email, user.ID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
