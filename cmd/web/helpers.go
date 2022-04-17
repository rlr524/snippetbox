package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name. If no entry exists
	// in the cache with the provided name, call the serverError helper method.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	// Execute the template set, passing in any dynamic data
	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}

// The serverError helper writes an error message and stack trace to the errorLog, then sends
// a generic 500 Internal Server Error response to the server
func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description to the user.
// This can be used to send responses like 400 "Bad Request" when there is a problem with the request the user sent.
func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, implement a notFound helper. This is a convenience wrapper around clientError which
// sends a 404 Not Found response to the user.
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
