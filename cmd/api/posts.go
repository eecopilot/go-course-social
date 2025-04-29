package main

import (
	"errors"
	"net/http"

	"github.com/eecopilot/go-course-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// var post store.Post // bad way
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	userId := 1

	// create a new post
	post := &store.Post{
		UserID:  int64(userId),
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")

	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			// writeJSONError(w, http.StatusInternalServerError, err.Error())
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		app.internalServerError(w, r, err)
		return
	}
}
