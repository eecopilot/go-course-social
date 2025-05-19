package main

import (
	"net/http"

	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"

	"github.com/eecopilot/go-course-social/internal/mailer"
	"github.com/eecopilot/go-course-social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}
type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		注册用户
//	@Description	注册一个新用户
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"用户信息"
//	@Success		201		{object}	UserWithToken		"用户创建成功"
//	@Failure		400		{object}	error				"请求错误"
//	@Failure		409		{object}	error				"用户已存在"
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// store user
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	// hash password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// token 生成
	// plainToken 发送给用户(email)
	plainToken := uuid.New().String()

	// 加密 token
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// 创建用户并发送邀请
	err := app.store.Users.CreateAndInvite(r.Context(), user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail, store.ErrDuplicateUsername:
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := &UserWithToken{
		User:  user,
		Token: plainToken,
	}

	isProdEnv := app.config.env == "production"
	activationURL := fmt.Sprintf("%s/confim/%s", app.config.frontendURL, plainToken)
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// 发送邮件
	err = app.mailer.Send(mailer.UserInvitationTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("send email failed", "error", err)
		// roll back user creation
		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.logger.Errorw("roll back user creation failed", "error", err)
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, r, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}
