-- 给 users 表添加 avatar 字段

BEGIN;

-- 添加 avatar 字段，类型为 TEXT（可以存储 URL 或文件路径）
-- 如果需要默认值，可以加上 DEFAULT 'default_avatar_url'
ALTER TABLE users ADD COLUMN avatar VARCHAR(255);

COMMIT;
