package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"oogway/first/snippetbox/internal/models"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLogger    *log.Logger
	infoLogger     *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	addr := flag.String("addr", ":4000", "Port the application runs on")
	dsn := flag.String("dsn", "web:Maverick2020!@/snippetbox?parseTime=true", "MySQL data source name")

	templateCache, err := newTemplateCache()

	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()

	db, err := openBD(*dsn)

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		infoLogger:     infoLog,
		errorLogger:    errorLog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	flag.Parse()

	infoLog.Printf("Starting server on %s", *addr)
	//err := http.ListenAndServe(*addr, mux)

	ser := &http.Server{

		Addr:     *addr,
		Handler:  app.routers(),
		ErrorLog: errorLog,
	}
	err = ser.ListenAndServe()
	errorLog.Fatal(err)
}

func openBD(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
