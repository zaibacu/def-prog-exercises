package main

import (
	"context"
	"database/sql"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	_ "embed"

	_ "github.com/glebarez/go-sqlite"
)

//go:embed notes.html
var notesTplSrc string

var notesTpl = template.Must(template.New("notes").Parse(notesTplSrc))

type note struct {
	id             int
	Title, Content string
}

func putNote(db *sql.DB, n note) error {
	_, err := db.Exec(`INSERT INTO notes(title, content) VALUES('` + n.Title + `', '` + n.Content + `')`)
	return err
}

func getNote(ctx context.Context, db *sql.DB, id int) (note, error) {
	rows, err := db.QueryContext(ctx, `SELECT * FROM notes WHERE id = `+strconv.Itoa(id))
	if err != nil {
		return note{}, err
	}
	if !rows.Next() {
		return note{}, errors.New("not found")
	}
	var n note
	if err := rows.Scan(&(n.id), &(n.Title), &(n.Content)); err != nil {
		return n, err
	}
	return n, rows.Err()
}

func getNotes(ctx context.Context, db *sql.DB) ([]note, error) {
	// Retrieve notes
	rows, err := db.QueryContext(ctx, `SELECT * FROM notes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []note
	for rows.Next() {
		// TODO use note instead
		var i int
		var t, c string
		if err := rows.Scan(&i, &t, &c); err != nil {
			return nil, err
		}
		notes = append(notes, note{i, t, c})
	}
	if err := rows.Err(); err != nil {
		return nil, rows.Err()
	}
	return notes, nil
}

func notes() http.Handler {

	db := must(sql.Open("sqlite", "./notes.db"))
	must(db.Exec(`CREATE TABLE IF NOT EXISTS notes(id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, content TEXT)`))
	if n, err := getNotes(context.Background(), db); err != nil || len(n) == 0 {
		if err := putNote(db, note{Title: "Salutations", Content: "Hello, World!"}); err != nil {
			panic(err)
		}
	}

	n := http.NewServeMux()

	// Home for the note page
	n.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		notes, err := getNotes(r.Context(), db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return
		}
		// Write the template with the notes
		if err = notesTpl.Execute(w, notes); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}

	})

	return n
}

func echo(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

func main() {
	sm := http.NewServeMux()
	sm.HandleFunc("GET /echo", echo)
	sm.Handle("/notes", notes())

	s := &http.Server{
		Addr:    "localhost:8080",
		Handler: sm,
	}
	log.Fatal(s.ListenAndServe())
}

func must[T any](t T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return t
}

// TODO authenticated section
// TODO handle note insertion
// TODO note deletion
// TODO styles
