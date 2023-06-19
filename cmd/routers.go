package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Chaque page commence avec chi.NewRouter()
func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Home page
	r.Get("/", app.home)

	// web page to retrieve data in json format
	// from server to web page
	r.Get("/jsonGraph", app.jsonData)

	// Source pages
	r.Get("/source/view/{id}", app.sourceView)
	r.Get("/source/create", app.sourceCreate)
	r.Post("/source/create", app.sourceCreatePost)
	r.Post("/source/delete/{id}", app.sourceDeletePost)
	r.Get("/source/update/{id}", app.sourceUpdate)
	r.Post("/source/update/{id}", app.sourceUpdatePost)

	// Info pages
	r.Get("/source/{sid}/info/view/{id}", app.infoView)
	r.Get("/source/{id}/info/create", app.infoCreate)
	r.Post("/source/{id}/info/create", app.infoCreatePost)
	r.Post("/source/{sid}/info/delete/{id}", app.infoDeletePost)
	r.Get("/source/{sid}/info/update/{id}", app.infoUpdate)
	r.Post("/source/{sid}/info/update/{id}", app.infoUpdatePost)

	// Static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return r
}
