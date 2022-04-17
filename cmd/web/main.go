package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rlr524/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	SnippetsModel *mysql.SnippetModel //SnippetsModel points to the SnippetModel struct that wraps the DB connection pool
}

func main() {
	// Command line flag for the port
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Command line flag for the MySQL DSN
	dsn := flag.String("dsn", "web:yukiorun!@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// Close the db connection pool before the main() function exits
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			errorLog.Fatal()
		}
	}(db)

	// Initialize an instance of Application containing logging dependencies
	app := &Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		SnippetsModel: &mysql.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Write messages using the infoLog and errorLog loggers, instead of the standard logger
	infoLog.Printf("Starting server on port %v", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
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
