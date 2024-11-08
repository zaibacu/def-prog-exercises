package main

import (
	"context"
	"io"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"

	"github.com/empijei/def-prog-exercises/app"
)

func main() {
	ctx := context.Background()
	sm := http.NewServeMux()
	sm.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, r.URL.Path+` doesn't exist on this server`)
	}))
	auth := app.Auth(ctx)
	sm.HandleFunc("/echo", app.Echo)
	sm.Handle("/auth/", auth)
	sm.Handle("/notes/", app.Notes(ctx))
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
