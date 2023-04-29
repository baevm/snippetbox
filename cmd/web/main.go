package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"snippetbox/internal/models"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type app struct {
	debug          bool
	errLogger      *log.Logger
	infoLogger     *log.Logger
	snippets       models.SnippetRepo
	users          models.UserRepo
	templaceCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

var config struct {
	addr  string
	dbDsn string
}

func main() {
	flag.StringVar(&config.addr, "addr", ":5000", "HTTP network address")
	flag.StringVar(&config.dbDsn, "dsn", "root:password@/snippetbox?parseTime=true", "MySQL connect name")
	debug := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(config.dbDsn)

	if err != nil {
		errLogger.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()

	if err != nil {
		errLogger.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := app{
		errLogger:      errLogger,
		infoLogger:     infoLogger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templaceCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		debug:          *debug,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	// custom config for server
	srv := &http.Server{
		Addr:         config.addr,
		ErrorLog:     errLogger,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLogger.Printf("Started listening on: %s", config.addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errLogger.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
