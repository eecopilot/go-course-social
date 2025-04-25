package main

import (
	"github.com/eecopilot/go-course-social/internal/env"
	"github.com/eecopilot/go-course-social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewStorage(nil)
	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		panic(err)
	}
}
