package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/eecopilot/go-course-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags" validate:"max=5"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// var post store.Post // bad way
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate the payload
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// 如果 Tags 为 nil，初始化为空数组
	if payload.Tags == nil {
		payload.Tags = []string{}
	}

	// create a new post
	post := &store.Post{
		UserID:  1,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, r, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r.Context())

	// 获取评论
	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err := app.jsonResponse(w, r, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r.Context())

	var payload UpdatePostPayload
	// readJSON 函数从 HTTP 请求体中读取 JSON 数据并解析到 payload 结构体中
	// 它会自动将请求体中的 JSON 字段映射到 UpdatePostPayload 结构体的对应字段
	// 例如，如果请求体是 {"title": "新标题", "content": "新内容"}
	// 则 payload.Title 将指向 "新标题"，payload.Content 将指向 "新内容"
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate the payload
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// update the post
	// 空指针判断
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, r, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// fetch post
func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId := chi.URLParam(r, "postId")
		ctx := r.Context()
		post, err := app.store.Posts.GetByID(ctx, postId)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(ctx context.Context) *store.Post {
	post, _ := ctx.Value(postCtx).(*store.Post)
	return post
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")

	ctx := r.Context()
	if err := app.store.Posts.Delete(ctx, postId); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, r, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
