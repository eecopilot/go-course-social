# go social course

# init

```
go mod init github.com/eecopilot/go-course-social
```

## chi

```
go get -u github.com/go-chi/chi/v5
```

## live reload

go install github.com/air-verse/air@latest

## env

echo export FOO=foo > .envrc
direnv allow .

## migrations

```sql
# -seq 告诉 migrate 使用顺序编号来代替时间戳（例如 000001_create_users.up.sql）
# -ext sql 表示你想要纯文本的 SQL 文件（即生成 .up.sql 和 .down.sql 文件）。
# -dir: 迁移文件存放目录
migrate create -seq -ext sql -dir ./cmd/migrate/migrations create_users
```

```sql
migrate -path=./cmd/migrate/migrations -database="$DATADASE_URL" up
```

## validator

go get github.com/go-playground/validator/v10
