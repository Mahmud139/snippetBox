package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Mahmud139/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: snippets,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: snippet,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
	//w.Write([]byte("Create a new snippet"))
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errs := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errs["title"] = "this field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errs["title"] =  "This field is too long (maximum is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errs["content"] = "this field cannot be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errs["expires"] = "this field cannot be blank"
	} else if expires != "365" && expires != "30" && expires != "7" {
		errs["expires"] = "This field is invalid"
	}

	if len(errs) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormData: r.PostForm,
			FormErrors: errs,
		})
		
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}
	
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}