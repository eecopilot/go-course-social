package store

import (
	"context"
	"database/sql"
	"time"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, string) (*Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, string) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetByID(context.Context, int64) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
		Activate(ctx context.Context, token string) error
		Delete(ctx context.Context, userID int64) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		Unfollow(ctx context.Context, followerID, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostgresPosts{db: db},
		Users:     &PostgresUsers{db: db},
		Comments:  &CommentsStores{db: db},
		Followers: &PostgresFollowers{db: db},
	}
}

func withTransaction(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	// start the transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// execute the function
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	// commit the transaction
	return tx.Commit()
}
