package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	// Use the alice package for middleware chain with the standard middleware used for every request
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// Create a file server which serves files out of the "./ui/static" directory. Note that
	// the path given to the http.Dir function is relative to the project directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static", http.StripPrefix("/static", fileServer))

	// Wrap the return statement in the recoverPanic and logRequest middleware, then pass the servemux as
	// the "next" parameter to the secureHeaders middleware. Because secureHeaders is just a function, and
	// the function returns a http.Handler, there is nothing else to do.
	return standardMiddleware.Then(mux)
}
