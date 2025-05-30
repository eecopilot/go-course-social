package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	app := newTestRequest(t)
	mux := app.mount()

	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		// check for 401 status code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		// if rr.Code != http.StatusUnauthorized {
		// 	t.Errorf("expected status code %d but got %d", http.StatusUnauthorized, rr.Code)
		// }

		// 使用 checkResponseCode 函数来检查响应状态码，代替if语句
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// add token to request header
		req.Header.Add("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)

		fmt.Println(rr.Body.String())
	})
}
