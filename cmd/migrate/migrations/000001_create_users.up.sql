
-- type User struct {
-- 	ID        int64  `json:"id"`
-- 	Username  string `json:"username"`
-- 	Email     string `json:"email"`
-- 	Password  string `json:"-"`
-- 	CreatedAt string `json:"created_at"`
-- }
CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email citext NOT NULL UNIQUE,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
