package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
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

	// Create an instance of a templateData struct holding the slice of snippets
	data := &templateData{Snippets: s}

	// Init a slice containing the paths to the show.page.gohtml file, plus the base layout and footer partial.
	// Note this files slice is different from the one used in the showSnippet handler.
	var files = []string{
		"./ui/html/home.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func navError(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("There's no page here, Madison."))
	if err != nil {
		log.Fatal("There was a problem with the NavError route", err)
	}
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

	// Create an instance of a templateData struct holding the snippet data
	data := &templateData{Snippet: s}

	// Init a slice containing the paths to the show.page.gohtml file, plus the base layout and footer partial.
	var files = []string{
		"./ui/html/show.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}

	// Parse the template files; pass in a spread operator to ParseFiles to iterate over the whole files slice
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Execute the template files. Note that snippet data (models.Snippet struct s) is passed as the final parameter
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
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
