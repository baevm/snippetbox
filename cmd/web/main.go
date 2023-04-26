package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type app struct {
	errLogger  *log.Logger
	infoLogger *log.Logger
}

var config struct {
	addr string
}

func main() {
	flag.StringVar(&config.addr, "addr", ":5000", "HTTP network address")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := app{
		errLogger:  errLogger,
		infoLogger: infoLogger,
	}

	mux := http.NewServeMux()

	// file server that serves ui/static folder
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("/ui/static/")})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/snippet/view/", app.SnippetView)
	mux.HandleFunc("/snippet/create", app.SnippetCreate)

	// use custom logger for server
	srv := &http.Server{
		Addr:     config.addr,
		ErrorLog: errLogger,
		Handler:  mux,
	}

	infoLogger.Printf("Started listening on: %s", config.addr)
	err := srv.ListenAndServe()
	errLogger.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
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
