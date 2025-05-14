package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

// password 是一个嵌套的结构体，用于存储密码的明文和哈希值
type Password struct {
	Text *string
	Hash []byte
}

func (p *Password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Hash = hash
	p.Text = &plainText
	return nil
}

type PostgresUsers struct {
	db *sql.DB
}

func (p *PostgresUsers) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username,password, email)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	err := tx.QueryRowContext(
		ctx, query,
		user.Username,
		user.Password.Hash,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (p *PostgresUsers) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var user User
	// Scan 是回填查询结果到query列出的变量中
	err := p.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (p *PostgresUsers) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTransaction(p.db, ctx, func(tx *sql.Tx) error {
		// 1. create the user
		if err := p.Create(ctx, tx, user); err != nil {
			return err
		}
		// 2. create the invitation
		err := p.createUserInvitation(ctx, tx, token, invitationExp, user.ID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (p *PostgresUsers) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	// try to insert the invitation
	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}
