package main

import (
	"net/http"

	"github.com/eecopilot/go-course-social/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		获取用户动态
//	@Description	获取用户动态
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"每页条数"	default(2)
//	@Param			offset	query		int		false	"偏移量"	default(0)
//	@Param			sort	query		string	false	"排序方式"	default("desc")
//	@Param			tags	query		string	false	"标签"
//	@Param			search	query		string	false	"搜索"
//	@Param			since	query		string	false	"开始时间"
//	@Param			until	query		string	false	"结束时间"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
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
