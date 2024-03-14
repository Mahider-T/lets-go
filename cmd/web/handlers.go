package main

import (
	"errors"
	"fmt"
	"net/http"
	"oogway/first/snippetbox/internal/models"
	"strconv"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		//http.NotFound(w, r)
		a.notFound(w)
		return
	}

	snippets, err := a.snippets.Latest()

	if err != nil {
		a.serverError(w, err)
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "+%v", snippet)
	}

	// files := []string{
	// "./ui/html/base.tmpl",
	// "./ui/html/partials/nav.tmpl",
	// "./ui/html/pages/home.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// log.Println(err.Error())
	//
	// a.errorLogger.Println(err.Error())
	// http.Error(w, "Internal server error", http.StatusInternalServerError)
	//
	// a.serverError(w, err)
	// return
	// }

	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	//log.Println(err.Error())
	//a.errorLogger.Println(err.Error())
	//http.Error(w, "Internal server error", http.StatusInternalServerError)

	// a.serverError(w, err)
	// }
}

func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		//http.NotFound(w, r)
		a.notFound(w)
		return
	}

	snp, err := a.snippets.Get(id)

	if err != nil {
		if errors.Is(models.ErrNoRecord, err) {
			a.notFound(w)
		} else {
			a.serverError(w, err)
		}

		return
	}
	// fmt.Fprintf(w, "Display a specific snippet with id ID %d ...", id)
	fmt.Fprintf(w, "%+v", snp)
}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n-Kobayashi Issa"
	expires := 7

	id, err := a.snippets.Insert(title, content, expires)

	if err != nil {
		a.serverError(w, err)
		return
	}

	// w.Write([]byte("Create a new snippet ..."))
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
