package main

import (
	"github.com/rlr524/snippetbox/handlers"
	"log"
	"net/http"
)

func main() {
	// Locally scoped servemux...never use the default servemux, i.e. http.HandleFunc
	// as it exposes the mux globally to all packages, including potentially malicious
	// In most cases will use a 3P mux like Chi or Gorilla and use it with a localally scoped var
	mux := http.NewServeMux()
	// Fixed path, think if it like /* (using a control flow to prevent it being a catch all in the handler)
	mux.HandleFunc("/", handlers.Home)
	// Subtree paths, only match when the path is exect
	mux.HandleFunc("/snippet", handlers.ShowSnippet)
	mux.HandleFunc("/snippet/create", handlers.CreateSnippet)
	mux.HandleFunc("/404", handlers.NavError)
	
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}