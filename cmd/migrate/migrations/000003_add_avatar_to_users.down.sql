-- 撤销添加 avatar 字段的操作

BEGIN;

-- 删除 avatar 字段
ALTER TABLE users DROP COLUMN avatar;

COMMIT;
