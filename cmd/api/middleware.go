package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// BasicAuthExplanation 解释HTTP中基本认证的工作原理
//
// 基本认证格式:
// 1. 客户端将用户名和密码用冒号组合: "username:password"
// 2. 这个字符串使用Base64编码
// 3. 编码后的字符串添加到Authorization头部，前缀为"Basic "
//
// 示例:
// - 用户名: admin
// - 密码: password2
// - 组合后: "admin:password2"
// - Base64编码: "YWRtaW46cGFzc3dvcmQy"
// - 最终头部: "Authorization: Basic YWRtaW46cGFzc3dvcmQy"
//
// 当服务器接收到这个头部时，它会:
// 1. 提取Base64编码的字符串
// 2. 将其解码回"username:password"
// 3. 分割字符串并验证凭据
//
// 注意: Base64是一种编码方式，而非加密。基本认证应该只在HTTPS下使用。

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//1. read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("no auth header provided"))
				return
			}
			//2. parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("no auth header provided"))
				return
			}
			//3. decode the base64 -> username:password
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid auth"))
			}
			//4. check the credentials
			creds := strings.Split(string(decoded), ":")
			if len(creds) != 2 {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid auth"))
				return
			}

			username := app.config.auth.basicAuth.username
			password := app.config.auth.basicAuth.password

			//5. check the credentials
			if username != creds[0] || password != creds[1] {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			//6. call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
