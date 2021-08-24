package main

import (
	"log"
	"net/http"
	"path/filepath"
)
const port = ":4000"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/snipped", ShowSnippet)
	mux.HandleFunc("/snippet/create", CreateSnippet)
	// Create a file server which serves files out of the "./ui/static" directory. Note that
	// the path given to the http.Dir function is relative to the project directory root.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	// Use the mux.Handle() function to register the file server as the handler for all URL paths
	// that start with "/static/". For matching paths, strip out the "/static" prefix
	// before the request reaches the file server.
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	log.Printf("Starting server on port %v", port)
	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}

// Custom file system to prohibit directory listing
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	
	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}