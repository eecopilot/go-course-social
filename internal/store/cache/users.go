package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eecopilot/go-course-social/internal/store"
	"github.com/go-redis/redis/v8"
)

type UserStore struct {
	rdb *redis.Client
}

const userExpTime = time.Minute

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	// s.rdb 是 redis 的客户端
	cacheKey := fmt.Sprintf("user-%v", userID)
	data, err := s.rdb.Get(ctx, cacheKey).Result()
	// redis.Nil 是一个错误类型，表示 key 不存在
	// 所以如果 key 不存在，我们返回 nil
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	user := store.User{}
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, data, userExpTime).Err()
}
