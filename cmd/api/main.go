package main

import (
	"log"

	"github.com/eecopilot/go-course-social/internal/db"
	"github.com/eecopilot/go-course-social/internal/env"
	"github.com/eecopilot/go-course-social/internal/store"
)

const version = "1.0.0"

func main() {
	cfg := config{
		addr:    env.GetString("ADDR", ":8080"),
		version: version,
		env:     env.GetString("ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres-db:5432"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 20),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 20),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "20s"),
		},
	}
	// db connection
	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panicf("failed to connect to db: %v", err)
	}

	defer db.Close()
	log.Println("connected to db")
	// store connection
	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}

	// 获取路由器
	mux := app.mount()

	// 启动HTTP服务器
	if err := app.run(mux); err != nil {
		panic(err)
	}
}
