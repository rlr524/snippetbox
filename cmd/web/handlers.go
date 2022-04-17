package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rlr524/snippetbox/pkg/models"
)

// The home function is defined as a method against *Application (a function receiver) (
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) // Use the notFound helper
		return
	}

	s, err := app.SnippetsModel.Latest()
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound helper
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a specific record based
	// on its ID. If no matching record is found, return a 404 Not Found response.
	s, err := app.SnippetsModel.Get(id)
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

// createSnippet function creates a new snippet #docs.md: createSnippet
func (app *Application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError helper
		return
	}
	// Dummy data for db
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nslowly!\n\n-Kobayashi Issa"
	expires := "7"

	// Pass the data to the SnippetModel.Insert() method, receiving the ID of the new record back.
	id, err := app.SnippetsModel.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Cache-Control", "public")
	_, err = w.Write([]byte("Create a new snippet..."))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
}
