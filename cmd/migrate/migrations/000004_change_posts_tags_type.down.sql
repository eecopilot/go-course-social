-- 1. 添加一个临时列来存储转换后的数据
ALTER TABLE posts ADD COLUMN tags_temp text[];

-- 2. 将 varchar(100) 转换为 text[]（将其作为数组的第一个元素）
UPDATE posts SET tags_temp = ARRAY[tags]::text[] WHERE tags IS NOT NULL;

-- 3. 删除原来的 tags 列
ALTER TABLE posts DROP COLUMN tags;

-- 4. 将临时列重命名为 tags
ALTER TABLE posts RENAME COLUMN tags_temp TO tags;

-- 5. 设置 NOT NULL 约束（如果需要）
ALTER TABLE posts ALTER COLUMN tags SET NOT NULL;
