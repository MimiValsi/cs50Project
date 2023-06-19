package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Web status are managed here

// serverError writes the error message
// and then sends 500 Internal Server Error to user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Print(trace)
}

// clientError send a specific status and describes
// to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound do the same thing but with 404 error
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Allocates memory so that a template can be rendered
// and checks if it exists before beeing sent
// to http.ResponseWriter
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	// Retrieves the appropriate template from cache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist",
			page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	// Executes the templates and than send to response bodyn
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	w.WriteHeader(status)

	buf.WriteTo(w)

}

// newTemplateData return a pointer to templateData
// no initialize and it's used by all functions in handlers.go file
// Make a better readability
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{}
}
