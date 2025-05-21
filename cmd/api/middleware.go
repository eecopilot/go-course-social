package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

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
				app.unauthorizedResponse(w, r, fmt.Errorf("no auth header provided"))
				return
			}
			//3. decode the base64 -> username:password
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedResponse(w, r, fmt.Errorf("invalid auth header"))
			}
			//4. check the credentials
			creds := strings.Split(string(decoded), ":")
			if len(creds) != 2 {
				app.unauthorizedResponse(w, r, fmt.Errorf("invalid auth header"))
				return
			}

			username := app.config.auth.basicAuth.username
			password := app.config.auth.basicAuth.password

			//5. check the credentials
			if username != creds[0] || password != creds[1] {
				app.unauthorizedResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			//6. call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
