package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

type Config struct {
	Port string
	StaticDir string
}

func main() {
	// Define command line flags for the port and static directory and parse them.
	cfg := new(Config)
	flag.StringVar(&cfg.Port, "port", ":4000", "HTTP port")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	// Create a file server which serves files out of the "./ui/static" directory. Note that
	// the path given to the http.Dir function is relative to the project directory root.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(cfg.StaticDir)})
	// Use the mux.Handle() function to register the file server as the handler for all URL paths
	// that start with "/static/". For matching paths, strip out the "/static" prefix
	// before the request reaches the file server.
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	log.Printf("Starting server on port %v", cfg.Port)
	err := http.ListenAndServe(cfg.Port, mux)
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