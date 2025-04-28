-- type Post struct {
-- 	ID        int64    `json:"id"`
-- 	UserID    int64    `json:"user_id"`
-- 	Title     string   `json:"title"`
-- 	Tags      []string `json:"tags"`
-- 	Content   string   `json:"content"`
-- 	CreatedAt string   `json:"created_at"`
-- 	UpdatedAt string   `json:"updated_at"`
-- }
CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id),
    title text NOT NULL,
    tags text[] NOT NULL,
    content text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
