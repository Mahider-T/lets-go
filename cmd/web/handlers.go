package main

import (
	"errors"
	"fmt"
	"net/http"
	"oogway/first/snippetbox/internal/models"
	"oogway/first/snippetbox/internal/validator"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
)

type userSignupForm struct {
	Name                string `form:"name"`
	Password            string `form:"password"`
	Email               string `form:"email"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

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
	// data.Flash = flash

	a.render(w, http.StatusOK, "view.tmpl", data)
}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := a.newTemplateModel(r)

	form := snippetCreateForm{
		Expires: 365,
	}
	data.Form = form
	a.render(w, http.StatusOK, "create.tmpl", data)
	// w.Write([]byte("Form to create a snippet zgoes here ..."))
}

func (a *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	err = a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field can not be empty"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field can not be more than 100 strings long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field can not be empty"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "Only a day, a week or a year are allowed in this field"
	}

	if len(form.FieldErrors) > 0 {
		data := a.newTemplateModel(r)
		data.Form = form
		a.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := a.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		a.serverError(w, err)
		return
	}

	a.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (a *application) userSignUp(w http.ResponseWriter, r *http.Request) {

	data := a.newTemplateModel(r)
	data.Form = userSignupForm{}
	a.render(w, http.StatusOK, "signup.tmpl", data)
	// fmt.Fprintln(w, "User sign up form goes here ...")
}

func (a *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {

	var form userSignupForm

	err := a.decodePostForm(r, &form)

	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := a.newTemplateModel(r)
		data.Form = form
		a.render(w, http.StatusNotAcceptable, "signup.tmpl", data)
		return
	}

	err = a.users.Insert(form.Name, form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email is already in use.")
			data := a.newTemplateModel(r)
			data.Form = form
			a.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		} else {
			a.serverError(w, err)
		}
		return
	}

	a.sessionManager.Put(r.Context(), "flash", "Your sign up was successful. Please login.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	// fmt.Fprintln(w, "Handler for the sing up ...")
}

func (a *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateModel(r)
	data.Form = userLoginForm{}

	a.render(w, http.StatusOK, "login.tmpl", data)
	// fmt.Fprintln(w, "Login form goes here ...")
}
func (a *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := a.decodePostForm(r, &form)

	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := a.newTemplateModel(r)
		data.Form = form
		a.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}
	id, err := a.users.Authenticate(form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := a.newTemplateModel(r)
			data.Form = form
			a.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
			return
		} else {
			a.serverError(w, err)
			return
		}
	}
	err = a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.serverError(w, err)
		return
	}

	a.sessionManager.Put(r.Context(), "authenticatedUserId", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	// fmt.Fprintln(w, "Login logic goes here...")
}

func (a *application) userLogout(w http.ResponseWriter, r *http.Request) {

	err := a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.serverError(w, err)
		return
	}
	a.sessionManager.Remove(r.Context(), "authenticatedUserId")

	a.sessionManager.Put(r.Context(), "flash", "You have been logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
	// fmt.Fprintln(w, "Logout a user ...")
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
