package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *settings) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.AuthTokenMW)

		// user routes
		mux.Post("/users", app.AllUsers)
		mux.Post("/users/save", app.EditUser)
		mux.Post("/users/get/{id}", app.GetUser)
		mux.Post("/users/delete", app.DeleteUser)
		mux.Post("/zuser/{id}", app.BootUser)
		// target routes
		mux.Post("/targets", app.AllTargets)
		mux.Post("/targets/save", app.EditTarget)
		mux.Post("/targets/get/{id}", app.GetTarget)
		mux.Post("/targets/delete", app.DeleteTarget)
		// config routes
		mux.Post("/configs", app.AllConfigs)
		mux.Post("/configs/get/{id}", app.GetConfig)
		mux.Post("/configs/delete", app.DeleteConfig)
		mux.Post("/configs/save", app.EditConfig)
		// console
		mux.Post("/send", app.SendCommand)
		mux.Post("/save/command", app.SaveCommand)
		mux.Post("/commands/get", app.GetUserSavedCommands)
		mux.Post("/responses", app.Last24responses)
		mux.Post("/responses/num/{num}", app.GetResponses)
		mux.Post("/responses/get/{id}", app.GetResponse)

	})

	mux.Route("/ipe", func(mux chi.Router) {
		mux.Use(app.AuthTokenMW)
		mux.Post("/auth", app.IpeAuth)
	})
	// non priv rts
	mux.Post("/vtk", app.ValidateToken)
	mux.Post("/users/logout", app.Logout)
	mux.Post("/users/login", app.Login)

	// static files
	fserver := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fserver))
	return mux
}
