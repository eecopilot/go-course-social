-- 恢复 tags 字段的 NOT NULL 约束
ALTER TABLE posts ALTER COLUMN tags SET NOT NULL; 