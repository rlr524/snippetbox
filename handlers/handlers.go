package handlers

import (
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := w.Write([]byte("Hello from Snippetbox"))
	if err != nil {
		log.Fatal("There was a problem with the home route")
	}
}

func NavError (w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("There's no page here, Madison."))
	if err != nil {
		log.Fatal("There was a problem with the NavError route")
	}
}

func ShowSnippet (w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Display a specific snippet"))
	if err != nil {
		log.Fatal("There was a problem with the ShowSnippet route")
	}
}

func CreateSnippet (w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Create a new snippet"))
	if err != nil {
		log.Fatal("There was a problem with the CreateSnippet route")
	}
}
