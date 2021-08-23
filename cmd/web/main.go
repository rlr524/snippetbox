package main

import (
	"log"
	"net/http"
)
const port = ":4000"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/snipped", ShowSnippet)
	mux.HandleFunc("/snippet/create", CreateSnippet)
	
	log.Printf("Starting server on port %v", port)
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}