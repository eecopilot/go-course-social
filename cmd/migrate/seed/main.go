package main

import (
	"log"

	"github.com/eecopilot/go-course-social/internal/db"
	"github.com/eecopilot/go-course-social/internal/env"
	"github.com/eecopilot/go-course-social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "")
	conn, err := db.NewDB(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewStorage(conn)
	db.Seed(store)
}
