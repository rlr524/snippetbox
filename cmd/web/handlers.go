package main

import (
	"errors"
	"fmt"
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

	for _, snippet := range s {
		fmt.Fprintf(w, "%v,\n", snippet)
	}

	//// Initialize a slice containing the paths of the two files. Note that the
	//// home.page.gohtml file must be the first file in the slice.
	//files := []string{
	//	"./ui/html/home.page.gohtml",
	//	"./ui/html/base.layout.gohtml",
	//	"./ui/html/footer.partial.gohtml",
	//}
	//
	//// Use the template.ParseFiles() function to read the files and store the templates in a
	//// template set. Note the ... here works very much like the spread operator in JavaScript as it
	//// unpacks the contents of the slice. If there is an error, log the detailed error message
	//// and use the http.Error() function to send a generic 500 Internal Server Error response to the user.
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err) // Use the serverError helper
	//	return
	//}
	//
	//// Then use the Execute() method from the html/template library on the template set to write the
	//// template content as the response body. The last parameter to Execute() represents any dynamic
	//// data that is passed in, which is currently nil.
	//err = ts.Execute(w, nil)
	//if err != nil {
	//	app.serverError(w, err) // Use the serverError helper
	//}
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

	// Write the snippet data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%v", s)

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
