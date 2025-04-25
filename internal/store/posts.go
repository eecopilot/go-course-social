package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64    `json:"id"`
	UserID    int64    `json:"user_id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostgresPosts struct {
	db *sql.DB
}

func (p *PostgresPosts) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (user_id, title, tags, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := p.db.QueryRowContext(
		ctx,
		query,
		post.UserID,
		post.Title,
		pq.Array(post.Tags),
		post.Content,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
