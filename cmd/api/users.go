package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/eecopilot/go-course-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userContextKey string

const userCtxKey userContextKey = "user"

// GetUserHandler	godoc
//
//	@Summary		获取用户信息
//	@Description	通过用户ID获取用户信息
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"用户ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, r, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

// followUserHandler	godoc
//
//	@Summary		关注用户
//	@Description	通过用户ID关注用户
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"用户ID"
//	@Success		204		{string}	string	"关注成功"
//	@Failure		400		{object}	error	"请求错误"
//	@Failure		404		{object}	error	"用户不存在"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r) // 关注者

	// TODO: 需要验证被关注者是否存在
	var payload FollowUser // 被关注者
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// 需要解决自己关注自己的问题
	if followerUser.ID == payload.UserID {
		app.conflictResponse(w, r, errors.New("cannot follow yourself"))
		return
	}

	if err := app.store.Followers.Follow(r.Context(), followerUser.ID, payload.UserID); err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicate):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, r, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// unfollowUserHandler godoc
//
//	@Summary		取消关注用户
//	@Description	通过用户ID取消关注用户
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"用户ID"
//	@Success		204		{string}	string	"取消关注成功"
//	@Failure		400		{object}	error	"请求错误"
//	@Failure		404		{object}	error	"用户不存在"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := getUserFromContext(r)

	// TODO: 需要验证被关注者是否存在
	var payload FollowUser // 被关注者
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.store.Followers.Unfollow(r.Context(), unfollowedUser.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, r, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// userContextMiddleware 从上下文中获取用户
func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getUserFromContext 从上下文中获取用户
func getUserFromContext(r *http.Request) *store.User {
	user, ok := r.Context().Value(userCtxKey).(*store.User)
	if !ok {
		return nil
	}
	return user
}
