package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eecopilot/go-course-social/internal/auth"
	"github.com/eecopilot/go-course-social/internal/store"
	"github.com/eecopilot/go-course-social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestRequest(t *testing.T) *application {
	t.Helper()

	// logger := zap.Must(zap.NewProduction()).Sugar()
	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected status code %d but got %d", expected, actual)
	}
}
