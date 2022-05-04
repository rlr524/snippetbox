package main

import "net/http"

func secureHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// This is a more common and simplified version of the secureHeaders function
// In my opinion, the above is more readable and makes it more clear that the fn function is a closure over
// the next handler.
//func secureHeaders(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("X-XSS-Protection", "1; mode=block")
//		w.Header().Set("X-Frame-Options", "deny")
//
//		next.ServeHTTP(w, r)
//	})
//}
