package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutesExist(t *testing.T) {
	paths := testApp.routes()
	chiRoutes := paths.(chi.Router)
	pathExists(t, chiRoutes, "/users/login")
}

func pathExists(t *testing.T, paths chi.Router, route string) {
	exists := false
	// walk the registered paths
	_ = chi.Walk(paths, func(method string, res string, h http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// if path is found
		if route == res {
			exists = true
		}
		return nil
	})
	if !exists {
		t.Errorf("path could not be found %v", route)
	}
}
