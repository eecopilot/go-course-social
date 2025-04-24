package main

import (
	"github.com/eecopilot/go-course-social/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}
	app := &application{
		config: cfg,
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		panic(err)
	}
}
