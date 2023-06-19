package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"CURATOR/database"

	// PostgreSQL driver
	"github.com/jackc/pgx/v4/pgxpool"
)

// Main struct, every struct within database folders connects here
type application struct {
	sources *database.Source
	infos   *database.Info

	templateCache map[string]*template.Template

	errorLog *log.Logger
	infoLog  *log.Logger

	DB *pgxpool.Pool
}

const (
	addr    = ":3005"
	dataURL = "postgres://web:pass@localhost:5432/cs50p"
)

func main() {

	// Ldate = Local data & Ltime = Local time
	infoLog := log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t",
		log.Ldate|log.Ltime|log.Lshortfile)

	// executes the comm function with DB
	db, err := openDB(dataURL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Fontion @ cmd/template.go
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		DB:      db,
		sources: &database.Source{},
		infos:   &database.Info{},

		templateCache: templateCache,

		infoLog:  infoLog,
		errorLog: errorLog,
	}

	// See routers.go
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

// Start communication with DB when needed
func openDB(dataURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	db, err := pgxpool.Connect(ctx, dataURL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
