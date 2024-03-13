package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLogger *log.Logger
	infoLogger  *log.Logger
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLogger:  infoLog,
		errorLogger: errorLog,
	}

	addr := flag.String("addr", ":4000", "Port the application runs on")

	flag.Parse()

	infoLog.Printf("Starting server on %s", *addr)
	//err := http.ListenAndServe(*addr, mux)

	ser := &http.Server{

		Addr:     *addr,
		Handler:  app.routers(),
		ErrorLog: errorLog,
	}
	err := ser.ListenAndServe()
	errorLog.Fatal(err)
}
