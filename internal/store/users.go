package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type PostgresUsers struct {
	db *sql.DB
}

func (p *PostgresUsers) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username,password, email )
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	err := p.db.QueryRowContext(
		ctx, query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
