-- 创建一个临时列
ALTER TABLE posts ADD COLUMN tags_temp text[];

-- 找出所有标签为NULL或空字符串的帖子，为它们设置一个空数组
UPDATE posts SET tags_temp = '{}' WHERE tags IS NULL OR tags = '';

-- 对于所有其他帖子，尝试将字符串转换为数组
-- 如果原标签不符合数组格式，则将其作为单个元素的数组
UPDATE posts 
SET tags_temp = 
    CASE
        WHEN tags LIKE '{%}' THEN tags::text[]
        ELSE ARRAY[tags]
    END
WHERE tags IS NOT NULL AND tags != '';

-- 删除原有的tags列
ALTER TABLE posts DROP COLUMN tags;

-- 重命名临时列
ALTER TABLE posts RENAME COLUMN tags_temp TO tags;

-- 设置NOT NULL约束（如果需要）
ALTER TABLE posts ALTER COLUMN tags SET NOT NULL; 