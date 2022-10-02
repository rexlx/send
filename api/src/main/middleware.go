package main

import (
	"fmt"
	"net/http"
)

func (app *settings) AuthTokenMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.models.Token.TestToken(r)
		if err != nil {
			data := jsonResponse{
				Error:   true,
				Message: fmt.Sprintf("found a dead one sir, %v", err),
			}
			_ = app.writeJSON(w, http.StatusUnauthorized, data)
			return
		}
		next.ServeHTTP(w, r)
	})
}
