package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// servemux "/" working like catch all
	// we need to check if path is really "/"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("snippet view"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// we need to check for method
	// because servemux doesnt support rest
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// http error shortcut
		// for WriteHeader, Write
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`"snippet":"textTestText"`))
}

func main() {
	srv := http.NewServeMux()

	srv.HandleFunc("/", home)
	srv.HandleFunc("/snippet/view/", snippetView)
	srv.HandleFunc("/snippet/create", snippetCreate)

	log.Println("Started listening on port: 5000")
	err := http.ListenAndServe(":5000", srv)
	log.Fatal(err)
}
