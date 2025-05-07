package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type PostgresFollowers struct {
	db *sql.DB
}

type Follower struct {
	UserID     int64  `json:"id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

func (p *PostgresFollowers) Follow(ctx context.Context, followerID, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := p.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return ErrDuplicate
			}
		}
	}
	return nil
}

func (p *PostgresFollowers) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `
		DELETE FROM followers WHERE user_id = $1 AND follower_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	_, err := p.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}
	return nil
}
