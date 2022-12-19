package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAllUsers(t *testing.T) {
	// need some mock rows
	var notRows = mockDB.NewRows([]string{
		"id", "email", "first_name", "last_name", "password", "user_active", "created_at", "updated_at", "has_token",
	})
	notRows.AddRow("42", "real@human.bot", "cappy", "bara", "hackerman", "1", time.Now(), time.Now(), "0")

	// expected a query that starts with `select`
	mockDB.ExpectQuery("select \\\\* ").WillReturnRows(notRows)

	// create a test recorder that satisfies the reqs for a response recorder
	rr := httptest.NewRecorder()
	// create req
	req, _ := http.NewRequest("POST", "/admin/users", nil)
	handler := http.HandlerFunc(testApp.AllUsers)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("bad code", rr.Code)
	}
}
