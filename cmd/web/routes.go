package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	// Use the alice package for middleware chain with the standard middleware used for every request
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	// Create a file server which serves files out of the "./ui/static" directory. Note that
	// the path given to the http.Dir function is relative to the project directory root.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	// Use the mux.Handle() function to register the file server as the handler for all URL paths
	// that start with "/static/". For matching paths, strip out the "/static" prefix
	// before the request reaches the file server.
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Wrap the return statement in the recoverPanic and logRequest middleware, then pass the servemux as
	// the "next" parameter to the secureHeaders middleware. Because secureHeaders is just a function, and
	// the function returns a http.Handler, there is nothing else to do.
	return standardMiddleware.Then(mux)
}
