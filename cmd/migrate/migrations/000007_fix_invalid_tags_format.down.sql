-- 这个回滚脚本不做任何实际改变，因为修复数据是单向操作
-- 我们不能准确地还原到原来的损坏状态

-- 如果需要回滚结构变更，可以执行：
-- ALTER TABLE posts ALTER COLUMN tags DROP NOT NULL;

-- 但不建议执行此操作，因为它可能会导致数据不一致 