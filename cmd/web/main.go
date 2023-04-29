package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"snippetbox/internal/models"
	"snippetbox/internal/templates"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	debug          bool
	errLogger      *log.Logger
	infoLogger     *log.Logger
	snippets       models.SnippetRepo
	users          models.UserRepo
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

var flags struct {
	addr  string
	dbDsn string
}

func main() {
	flag.StringVar(&flags.addr, "addr", ":5000", "HTTP network address")
	flag.StringVar(&flags.dbDsn, "dsn", "root:password@/snippetbox?parseTime=true", "MySQL connect name")
	debug := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(flags.dbDsn)

	if err != nil {
		errLogger.Fatal(err)
	}

	defer db.Close()

	templateCache, err := templates.NewTemplateCache()

	if err != nil {
		errLogger.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := App{
		errLogger:      errLogger,
		infoLogger:     infoLogger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
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
		Addr:         flags.addr,
		ErrorLog:     errLogger,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLogger.Printf("Started listening on: %s", flags.addr)
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
