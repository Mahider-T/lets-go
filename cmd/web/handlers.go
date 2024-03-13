package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		//http.NotFound(w, r)
		a.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		//log.Println(err.Error())

		//a.errorLogger.Println(err.Error())
		//http.Error(w, "Internal server error", http.StatusInternalServerError)

		a.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		//log.Println(err.Error())
		//a.errorLogger.Println(err.Error())
		//http.Error(w, "Internal server error", http.StatusInternalServerError)

		a.serverError(w, err)
	}
}

func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		//http.NotFound(w, r)
		a.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with id ID %d ...", id)
}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet ..."))
}
