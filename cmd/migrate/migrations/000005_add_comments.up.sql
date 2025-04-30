CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    user_id bigserial NOT NULL REFERENCES users(id),
    post_id bigserial NOT NULL REFERENCES posts(id),
    content TEXT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
