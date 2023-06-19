package main

import (
	"html/template"
	"path/filepath"

	"time"

	"CURATOR/database"
)

// Template struct that generate and analyse data
// to .html.gotpl files.
type templateData struct {
	Source  *database.Source
	Sources []*database.Source

	Info  *database.Info
	Infos []*database.Info

	JSource []byte

	Form any
}

// @ tables sources et infos, columns "Created" and "Updated"
// have timestamp (UTC)
// SELECT NOW()::timestamp;
// 2023-02-10 19:28:53.116296
// |________| needed
func humanDate(t time.Time) string {
	return t.Format("02/01/2006")
}

// template.FuncMap is stocked in a global variable
// so it's easier to used it with humanDate function
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Create a slice with all paths
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		/// extracts the path file name
		name := filepath.Base(page)

		// Create a new empty template
		ts, err := template.New(name).Funcs(functions).
			ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
