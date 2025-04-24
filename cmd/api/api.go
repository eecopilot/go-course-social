package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Group
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})

	return r
}

func (app *application) run(mux *chi.Mux) error {
	// consts for server settings
	const WTime = 20 * time.Second
	const RTime = 10 * time.Second
	const ITime = time.Minute

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: WTime, // max time to write response to the client
		ReadTimeout:  RTime, // max time to read request from the client
		IdleTimeout:  ITime, // max time for connections using TCP Keep-Alive
	}
	log.Printf("Starting server on http://localhost%s", app.config.addr)
	return srv.ListenAndServe()
}
