package main

import (
	"net/http"

	"github.com/eecopilot/go-course-social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// user := getUserFromContext(r)
	// log.Println(user)
	fq := store.PaginatedFeedQuery{
		Limit:  2,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
		Search: "",
		Since:  "",
		Until:  "",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// validation
	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()
	// TODO: 从上下文中获取用户
	feeds, err := app.store.Posts.GetUserFeed(ctx, 1, fq)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, r, http.StatusOK, feeds); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
