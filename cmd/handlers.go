package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"CURATOR/database"
	"CURATOR/internal/validator"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (app *application) dbConn(ctx context.Context) *pgxpool.Conn {
	conn, err := app.DB.Acquire(ctx)
	if err != nil {
		app.errorLog.Println("Unable to connect to DB.")
	}

	return conn
}

//
// Home
//

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	// MenuSource func @ database/sources.go
	sources, err := app.sources.MenuSource(conn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// A copy of the data send to web page but in json format
	// the JS function that generate the graphs can work proprely
	jData, err := json.Marshal(sources)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// newTemplateData @ cmd/templates.go
	data := app.newTemplateData(r)
	data.Sources = sources
	data.JSource = jData

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) jsonData(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	sources, err := app.sources.MenuSource(conn)
	if err != nil {
		app.serverError(w, err)
	}

	jsonGraph, err := json.Marshal(sources)
	if err != nil {
		app.serverError(w, err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonGraph)
}

//
// Sources Handlers
//

type sourceCreateForm struct {
	Name string

	validator.Validator
}

// Generate source view with a table of all infos within
func (app *application) sourceView(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Call database/sources.go function
	// Fetch 'id' from URL create before
	source, err := app.sources.SourceGet(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	// Call database/infos.go function
	// with source id
	info, err := app.infos.InfoList(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	data := app.newTemplateData(r)
	data.Infos = info
	data.Source = source

	app.render(w, http.StatusOK, "sourceView.tmpl.html", data)

}

// Source creation page.
func (app *application) sourceCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Form = sourceCreateForm{}

	app.render(w, http.StatusOK, "sourceCreate.tmpl.html", data)
}

func (app *application) sourceCreatePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	// parseForm fetch variable from URL
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := sourceCreateForm{
		Name: r.PostForm.Get("name"),
	}

	// No empty field helper
	emptyField := "Cannot be empty"

	form.CheckField(validator.NotBlank(form.Name),
		"name", emptyField)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity,
			"sourceCreate.tmpl.html", data)
		return
	}

	// if no error, than data it sent to DB
	id, err := app.sources.SourceInsert(form.Name, conn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", id),
		http.StatusSeeOther)
}

// Fetch source id from URL and send delete command to PSQL
func (app *application) sourceDeletePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	key := chi.URLParam(r, "id")

	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.sources.SourceDelete(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Fetch data from source id and save modifications
// made by user and send them to PSQL
func (app *application) sourceUpdate(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	source, err := app.sources.SourceGet(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Source = source

	app.render(w, http.StatusOK, "sourceUpdate.tmpl.html", data)
}

func (app *application) sourceUpdatePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil {
		app.notFound(w)
		return
	}

	form := sourceCreateForm{
		Name: r.PostForm.Get("name"),
	}

	emptyField := "Ce champ ne doit pas être vide"

	form.CheckField(validator.NotBlank(form.Name),
		"name", emptyField)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity,
			"sourceUpdate.tmpl.html", data)
		return
	}

	app.sources.Name = form.Name

	err = app.sources.SourceUpdate(id, conn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", id),
		http.StatusSeeOther)

}

//
// Infos Handlers
//

type infoCreateForm struct {
	ID       int
	Agent    string
	Material string
	Priority string
	Detail   string
	Created  string
	Updated  string
	Status   string
	Estimate string

	validator.Validator
}

// Same thing as source. Fetch source id so it can be sent
// to source_id (FK)
func (app *application) infoCreate(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	// Fetch source id from URL
	key := chi.URLParam(r, "id")

	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	source, err := app.sources.SourceGet(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Form = infoCreateForm{}
	data.Source = source

	app.render(w, http.StatusOK, "infoCreate.tmpl.html", data)
}

func (app *application) infoCreatePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	key := chi.URLParam(r, "id")

	sID, err := strconv.Atoi(key)
	if err != nil || sID < 1 {
		app.notFound(w)
		return
	}

	// Les données récupérés depuis la page HTML sont envoyées
	// vers la BD
	form := infoCreateForm{
		Agent:    r.PostForm.Get("agent"),
		Material: r.PostForm.Get("material"),
		Detail:   r.PostForm.Get("detail"),
		Priority: r.PostForm.Get("priority"),
		Estimate: r.PostForm.Get("estimate"),
		Status:   r.PostForm.Get("status"),
	}

	// These can't be empty
	// Below ensures that the user is alerted
	emptyField := "Cannot be empty"

	form.CheckField(validator.NotBlank(form.Agent),
		"agent", emptyField)
	form.CheckField(validator.NotBlank(form.Material),
		"material", emptyField)
	form.CheckField(validator.NotBlank(form.Detail),
		"detail", emptyField)
	form.CheckField(validator.NotBlank(form.Priority),
		"priority", emptyField)
	form.CheckField(validator.NotBlank(form.Status),
		"status", emptyField)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity,
			"infoCreate.tmpl.html", data)
		return
	}

	app.infos.Agent = form.Agent
	app.infos.Material = form.Material
	app.infos.Detail = form.Detail
	app.infos.Estimate = form.Estimate
	app.infos.Status = form.Status
	app.infos.Priority, err = strconv.Atoi(form.Priority)
	if err != nil {
		app.notFound(w)
		return
	}

	_, err = app.infos.Insert(sID, conn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", sID),
		http.StatusSeeOther)
}

// Show detailed data from info
func (app *application) infoView(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	iKey := chi.URLParam(r, "id")

	id, err := strconv.Atoi(iKey)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	info, err := app.infos.InfoGet(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Info = info

	app.render(w, http.StatusOK, "infoView.tmpl.html", data)
}

// delete info
func (app *application) infoDeletePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	sKey := chi.URLParam(r, "sid")
	iKey := chi.URLParam(r, "id")

	id, err := strconv.Atoi(iKey)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	sID, err := strconv.Atoi(sKey)
	if err != nil || sID < 1 {
		app.notFound(w)
		return
	}

	err = app.infos.InfoDelete(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", sID),
		http.StatusSeeOther)

}

// Updates info. Same behavior as sourceUpdate
func (app *application) infoUpdate(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// id ~> Info id
	info, err := app.infos.InfoGet(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Info = info

	app.render(w, http.StatusOK, "infoUpdate.tmpl.html", data)
}

func (app *application) infoUpdatePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	sKey := chi.URLParam(r, "sid")
	sID, err := strconv.Atoi(sKey)
	if err != nil || sID < 1 {
		app.notFound(w)
		return
	}

	iKey := chi.URLParam(r, "id")
	iID, err := strconv.Atoi(iKey)
	if err != nil || iID < 1 {
		app.notFound(w)
		return
	}

	form := infoCreateForm{
		Agent:    r.PostForm.Get("agent"),
		Material: r.PostForm.Get("material"),
		Detail:   r.PostForm.Get("detail"),
		Priority: r.PostForm.Get("priority"),
		Estimate: r.PostForm.Get("estimate"),
		Status:   r.PostForm.Get("status"),
	}

	app.infos.Agent = form.Agent
	app.infos.Material = form.Material
	app.infos.Detail = form.Detail
	app.infos.Estimate = form.Estimate
	app.infos.Status = form.Status
	app.infos.Priority, err = strconv.Atoi(form.Priority)
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.infos.InfoUpdate(iID, conn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/%d/info/view/%d",
		sID, iID), http.StatusSeeOther)
}
