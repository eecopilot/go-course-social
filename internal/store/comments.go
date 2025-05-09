package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type CommentsStores struct {
	db *sql.DB
}

func (s *CommentsStores) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
		SELECT c.id, c.user_id, c.post_id, c.content, c.created_at, u.username
		FROM comments as c
		LEFT JOIN users as u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 创建一个切片来存储评论
	comments := []Comment{}
	// 遍历结果集
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID, &c.UserID, &c.PostID, &c.Content, &c.CreatedAt, &c.User.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (s *CommentsStores) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments ( post_id, user_id, content)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	// 这里PostID和UserID的顺序不能颠倒，要跟query中的顺序一致 post_id, user_id, content
	// Scan是做什么的
	// Scan将查询结果的列映射到结构体的字段
	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}
