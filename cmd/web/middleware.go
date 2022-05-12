package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always run in the event of a panic as Go unwinds the stack).
		defer func() {
			// Use the built-in recover function to check if there has been a panic or not. If there has:
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")

				// Call the app.ServerError helper to return a 500 response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// This is a more common and simplified version of the secureHeaders function (and applicable to the other
// middleware functions as well). In my opinion, the above is more readable and makes it more clear that
// the fn function is a closure over the next handler.
//func secureHeaders(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("X-XSS-Protection", "1; mode=block")
//		w.Header().Set("X-Frame-Options", "deny")
//
//		next.ServeHTTP(w, r)
//	})
//}
