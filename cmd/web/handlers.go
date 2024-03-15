package main

import (
	"errors"
	"fmt"
	"net/http"
	"oogway/first/snippetbox/internal/models"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := a.snippets.Latest()

	if err != nil {
		a.serverError(w, err)
		return
	}

	data := a.newTemplateModel(r)
	data.Snippets = snippets

	a.render(w, http.StatusOK, "home.tmpl", data)
}

func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext((r.Context()))

	id, err := strconv.Atoi(params.ByName("id"))
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

	data := a.newTemplateModel(r)
	data.Snippet = snp

	a.render(w, http.StatusOK, "view.tmpl", data)
}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Form to create a snippet goes here ..."))
}

func (a *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n-Kobayashi Issa"
	expires := 7

	id, err := a.snippets.Insert(title, content, expires)

	if err != nil {
		a.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
