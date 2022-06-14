package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/rlr524/snippetbox/pkg/models"
)

// The home function is defined as a method against *Application (a function receiver) (
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) // Use the notFound helper
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the render helper.
	app.render(w, r, "home.page.gohtml", &templateData{
		Snippets: s,
	})
}

func (app *Application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound helper
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a specific record based
	// on its ID. If no matching record is found, return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the render helper.
	app.render(w, r, "show.page.gohtml", &templateData{
		Snippet: s,
	})
}

// createSnippetForm function is a handler for presenting to form used to create a new snippet
func (app *Application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.gohtml", nil)
}

// createSnippet function creates a new snippet #docs.md: createSnippet
func (app *Application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Call r.ParseForm which adds any data in POST request bodies to the r.PostForm map. This also works in the
	// same way for PUT and PATCH requests. If there are any errors, us the app.ClientError helper to send a 400.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method tp retrieve the relevant data fields from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	// Hold any validation errors
	errorsMap := make(map[string]string)

	// Check that the title field is not blank and is under 100 char. If it fails either check, add a message
	// to the errors map using the field name as the key. Using the RuneCountInString method here because we
	// want to count the characters in the string, not the number of bytes, which is what we'd get if we used
	// the len() function on a string. See https://go.dev/play/p/DETcUgcQgv3 for a Go playground that
	// demonstrates the difference between len() and RuneCountInString().
	if strings.TrimSpace(title) == "" {
		errorsMap["title"] = "The title field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errorsMap["title"] = "The title field is too long (maximum is 100 characters)"
	}

	// Check that the content field isn't blank
	if strings.TrimSpace(content) == "" {
		errorsMap["content"] = "The content field cannot be blank"
	}

	// Check that the expires field isn't blank and matches one of the permitted values ("1", "7" or "365)
	if strings.TrimSpace(expires) == "" {
		errorsMap["expires"] = "The expires field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errorsMap["expires"] = "The expires field contains an invalid value, it must be 1, 7 or 365"
	}

	// If there are any validation errors, redisplay the create.page.gohtml template, passing in the
	// validation errors and the previously submitted r.PostForm data.
	if len(errorsMap) > 0 {
		app.render(w, r, "create.page.gohtml", &templateData{
			FormErrors: errorsMap,
			FormData:   r.PostForm,
		})
		return
	}

	// Create a new snippet in the db using the form data
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Cache-Control", "public")
}
