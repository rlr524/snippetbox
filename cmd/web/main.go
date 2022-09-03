package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
	"github.com/rlr524/snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel //SnippetsModel points to the SnippetModel struct that wraps the DB connection pool
	templateCache map[string]*template.Template
}

func main() {
	_ = os.Setenv("environment", "development")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbPass := os.Getenv("DB_PASS")
	sessionSecret := os.Getenv("SESSION_SECRET")
	// Command line flag for the port
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Command line flag for the MySQL DSN currently located on local Docker container
	// TODO: Encrypt the password
	dsn := flag.String("dsn", fmt.Sprintf("web:%s@tcp(127.0.0.1:3306)/snippetbox?parseTime=true", dbPass),
		"MySQL data source name")
	secret := flag.String("secret", sessionSecret, "Secret key")
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

	// Initialize a new template cache.
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	// Initialize an instance of Application containing logging dependencies, models and cache
	app := &Application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Write messages using the infoLog and errorLog loggers, instead of the standard logger
	infoLog.Printf("Starting server on port %v", *addr)
	err = srv.ListenAndServeTLS("./security/cert.pem", "./security/key.pem")
	errorLog.Fatal(err)
}

// The openDB function returns a pool of connections for the MySQL DB
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
