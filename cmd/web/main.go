package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	Port      string
	StaticDir string
}

func main() {
	// Define command line flags for the port and static directory and parse them.
	cfg := new(Config)
	flag.StringVar(&cfg.Port, "port", ":4000", "HTTP port")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	// Use log.New() to create a logger for writing information messages. This takes three parameters:
	// 1. The destination to write the logs to (os.Stdout) 2. A string prefix for the message,
	// followed by a tab and 3. Flags to to indicate what additional information to include (local date and time).
	// Note that the the date and time are joined using a bitwise OR operator (|).
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	// Create a logger for writing error messages in the same way, but use stderr as the destination and use the
	// log.Llongfile flag to include the relevant path, file name and line number
	errorLog := log.New(os.Stderr, "ERR:\t", log.Ldate|log.Ltime|log.Llongfile)

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

	// Initialize a new http.Server struct. Set the Addr and Handler fields so that the server
	// uses the same network address as in the Config struct and the routes and set the ErrorLog field so
	// that the server uses the custom errorLog logger in the event of any problems.
	srv := &http.Server{
		Addr:     cfg.Port,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Write messages using the infoLog and errorLog loggers, instead of the standard logger
	infoLog.Printf("Starting server on port %v", cfg.Port)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
	/** log.Printf("Starting server on port %v", cfg.Port)
	err := http.ListenAndServe(cfg.Port, mux)
	log.Fatal(err) */
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
