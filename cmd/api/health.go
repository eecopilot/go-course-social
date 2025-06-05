package main

import (
	"net/http"
	"time"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "ok", "env": app.config.env, "version": app.config.version}

	// 模拟长时间响应
	time.Sleep(5 * time.Second)
	if err := app.jsonResponse(w, r, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
