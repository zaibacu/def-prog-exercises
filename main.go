package main

import (
	"io"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"

	"github.com/empijei/def-prog-exercises/app"
)

func main() {
	sm := http.NewServeMux()
	sm.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, r.URL.Path+` doesn't exist on this server`)
	}))
	sm.HandleFunc("/echo", app.Echo)
	// sm.HandleFunc("/login/", app.Login())
	sm.Handle("/notes/", app.Notes())
	addr := "localhost:8080"
	s := &http.Server{
		Addr:    addr,
		Handler: sm,
	}
	log.Println("Ready to accept connections on " + addr)
	log.Fatal(s.ListenAndServe())
}

// TODO authenticated section
// TODO note deletion
// TODO styles
