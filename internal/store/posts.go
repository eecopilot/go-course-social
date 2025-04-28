package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrNotFound = errors.New("resource not found")
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

func (p *PostgresPosts) GetByID(ctx context.Context, postId string) (*Post, error) {
	query := `
		SELECT id, user_id, title, tags, content, created_at, updated_at
		FROM posts
		WHERE id = $1
	`
	var post Post
	// Scan 方法将查询结果的各列值扫描到对应的变量中
	// 这里将查询返回的 id, user_id, title, tags, content, created_at, updated_at 列的值
	// 分别赋值给 post 结构体的对应字段
	// 如果查询没有返回行或发生其他错误，err 将不为 nil
	err := p.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,             // 将 id 列的值赋给 post.ID
		&post.UserID,         // 将 user_id 列的值赋给 post.UserID
		&post.Title,          // 将 title 列的值赋给 post.Title
		pq.Array(&post.Tags), // 使用 pq.Array 处理 PostgreSQL 数组类型
		&post.Content,        // 将 content 列的值赋给 post.Content
		&post.CreatedAt,      // 将 created_at 列的值赋给 post.CreatedAt
		&post.UpdatedAt,      // 将 updated_at 列的值赋给 post.UpdatedAt
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}
