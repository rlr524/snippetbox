package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	// Use the alice package for middleware chain with the standard middleware used for every request
	middleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, app.session.Enable)

	r := chi.NewRouter()

	r.Get("/", app.home)
	r.Route("/snippet", func(r chi.Router) {
		r.Get("/create", app.createSnippetForm)
		r.Post("/create", app.createSnippet)
		r.Get("/{id:[0-9]+}", app.showSnippet)
	})

	// Create a file server which serves files out of the "./ui/static" directory. Note that
	// the path given to the http.Dir function is relative to the project directory root.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	// Use the mux.Handle() function to register the file server as the handler for all URL paths
	// that start with "/static/". For matching paths, strip out the "/static" prefix
	// before the request reaches the file server.
	r.Handle("/static", http.NotFoundHandler())
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Wrap the return statement in the recoverPanic and logRequest middleware, then pass the servemux as
	// the "next" parameter to the secureHeaders middleware. Because secureHeaders is just a function, and
	// the function returns a http.Handler, there is nothing else to do.
	return middleware.Then(r)
}
