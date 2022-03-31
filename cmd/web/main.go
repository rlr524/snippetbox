package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Llongfile)

	// Initialize an instance of Application containing logging dependencies
	app := &Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Write messages using the infoLog and errorLog loggers, instead of the standard logger
	infoLog.Printf("Starting server on port %v", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
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
