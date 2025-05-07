package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var (
	ErrNotFound  = errors.New("resource not found")
	ErrDuplicate = errors.New("duplicate resource")
	QueryTimeout = 5 * time.Second
)

type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	// 查询数据库
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
		SELECT id, user_id, title, tags, content, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
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
		&post.Version,        // 将 version 列的值赋给 post.Version
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

// 加入version字段，防止并发更新
func (p *PostgresPosts) Update(ctx context.Context, payload *Post) error {
	query := `
		UPDATE posts SET title = $1, content = $2, version = version + 1 WHERE id = $3 AND version = $4 RETURNING version
	`
	// QueryRowContext 方法执行一个预编译的 SQL 语句，并返回一个 sql.Row 对象
	// 如果查询没有返回行或发生其他错误，err 将不为 nil

	// Scan方法用于将查询结果的列值扫描到指定的变量中
	// 这里执行UPDATE语句后，使用RETURNING子句返回了更新后的version值
	// Scan将这个返回的version值赋给payload.Version变量
	// 如果没有找到匹配的行(版本不匹配或ID不存在)，Scan会返回sql.ErrNoRows错误
	err := p.db.QueryRowContext(ctx, query, payload.Title, payload.Content, payload.ID, payload.Version).Scan(&payload.Version)
	// _, err := p.db.ExecContext(ctx, query, payload.Title, payload.Content, payload.ID, payload.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (p *PostgresPosts) Delete(ctx context.Context, postId string) error {
	query := `
		DELETE FROM posts WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	// ExecContext 方法执行一个预编译的 SQL 语句，并返回一个 sql.Result 对象
	res, err := p.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}
	// 确定受影响的行数
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	// 如果受影响的行数为 0，则返回 ErrNotFound
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
