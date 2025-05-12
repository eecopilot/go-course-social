package main

import (
	"github.com/eecopilot/go-course-social/internal/db"
	"github.com/eecopilot/go-course-social/internal/env"
	"github.com/eecopilot/go-course-social/internal/store"
	"go.uber.org/zap"
)

const version = "1.0.0"

//	@title			Swagger Social API
//	@description	This is a sample server Social API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				请在请求头中添加Authorization字段，值为Bearer + 空格 + token

func main() {
	cfg := config{
		addr:    env.GetString("ADDR", ":8080"),
		version: version,
		env:     env.GetString("ENV", "development"),
		apiUrl:  env.GetString("API_URL", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres-db:5432"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 20),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 20),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "20s"),
		},
	}
	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// db connection
	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal("failed to connect to db: %v", err)
	}

	defer db.Close()
	logger.Info("connected to db")
	// store connection
	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	// 获取路由器
	mux := app.mount()

	// 启动HTTP服务器
	if err := app.run(mux); err != nil {
		logger.Fatal("failed to start server: %v", err)
	}
}
