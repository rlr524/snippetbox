package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rlr524/snippetbox/pkg/forms"
	"github.com/rlr524/snippetbox/pkg/models"
	"net/http"
	"strconv"
)

// The home function is defined as a method against *Application (a function receiver) (
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
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
	app.render(w, r, "create.page.gohtml", &templateData{
		Form: forms.New(nil),
	})
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

	// The forms.Form struct contains the POSTed data from the form, then uses the validation methods to check content.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// If the form isn't valid, redisplay the template passing in the form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "create.page.gohtml", &templateData{
			Form: form,
		})
		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded in the form.Form struct,
	// use the Get() method to retrieve the validated value for a particular form field.
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value and the corresponding key to the session data.
	// Note that if there is no existing session for the current user (or their session has expired) then a
	// new empty session for them will be automatically created by the session middleware.
	app.session.Put(r, "toast", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Cache-Control", "public")
}

func (app *Application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.gohtml", &templateData{
		Form: forms.New(nil),
	})
}

func (app *Application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using the form helper
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	// If there are any errors, redisplay the signup form
	if !form.Valid() {
		app.render(w, r, "signup.page.gohtml", &templateData{Form: form})
		return
	}

	// Try to create a new user record in the db. If the email already exists, add an error
	// message to the form and redisplay it.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.FormErrors.Add("email", "Email address is already in use")
			app.render(w, r, "signup.page.gohtml", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Otherwise add a confirmation toast message to the session confirming that the signup worked
	// and asking the user to now log in.
	app.session.Put(r, "flash", "Your signup was successful, please log in.")

	// And redirect to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", &templateData{
		Form: forms.New(nil),
	})
}

func (app *Application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check if credentials are valid. If not, add a generic error message to the form failures
	// map and re-display the login page.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.FormErrors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.gohtml", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Add the ID of the current user to the session, so they are now "logged in".
	app.session.Put(r, "authenticatedUserID", id)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *Application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
