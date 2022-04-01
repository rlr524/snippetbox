package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// The home function is defined as a method against *Application (a function receiver) (
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) // Use the notFound helper
		return
	}

	// Initialize a slice containing the paths of the two files. Note that the
	// home.page.gohtml file must be the first file in the slice.
	files := []string{
		"./ui/html/home.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}

	// Use the template.ParseFiles() function to read the files and store the templates in a
	// template set. Note the ... here works very much like the spread operator in JavaScript as it
	// unpacks the contents of the slice. If there is an error, log the detailed error message
	// and use the http.Error() function to send a generic 500 Internal Server Error response to the user.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Use the serverError helper
		return
	}

	// Then use the Execute() method from the html/template library on the template set to write the
	// template content as the response body. The last parameter to Execute() represents any dynamic
	// data that is passed in, which is currently nil.
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err) // Use the serverError helper
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
	_, e := fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	if e != nil {
		log.Fatal("There was a problem with the ShowSnippet route", e)
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
	id, err := app.Snippets.Insert(title, content, expires)
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
