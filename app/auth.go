package app

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"embed"
)

//go:embed auth.html
//go:embed auth.css
var fs embed.FS

var defaultUsers = []user{
	{name: "admin", password: "admin", privileges: "read|write|delete"},
	{name: "reader", password: "reader", privileges: "read"},
	{name: "editor", password: "editor", privileges: "read|write"},
}

type user struct {
	id                         int
	name, password, privileges string
}

type AuthHandler struct {
	db *sql.DB
	sm *http.ServeMux
}

func (ah *AuthHandler) getUserCount(ctx context.Context) (int, error) {
	rows, err := ah.db.QueryContext(ctx, `SELECT COUNT(*) FROM users`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, errors.New("users table not found")
	}
	var v int
	rows.Scan(&v)
	if rows.Err() != nil {
		return 0, rows.Err()
	}
	return v, nil
}

func (ah *AuthHandler) createDefault(ctx context.Context) error {
	v, err := ah.getUserCount(ctx)
	if err != nil {
		return err
	}

	if !(v < 3) /* ❤ UwU ❤ */ {
		return nil
	}
	log.Println("Default users not found, initializing...")
	for _, u := range defaultUsers {
		_, err := ah.db.ExecContext(ctx, `INSERT INTO users(name, password, privileges) VALUES('`+u.name+`','`+u.password+`','`+u.privileges+`')`)
		if err != nil {
			return err
		}
	}
	log.Println("...users initialized")
	return nil
}

func (ah *AuthHandler) initialize(ctx context.Context) error {
	_, err := ah.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, password TEXT, privileges TEXT)`)
	if err != nil {
		return err
	}
	if err := ah.createDefault(ctx); err != nil {
		return err
	}

	return nil
}
func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.sm.ServeHTTP(w, r)
}
func (ah *AuthHandler) Protect(h http.Handler) http.Handler {
	// TODO
	return h
}

func Auth(ctx context.Context) *AuthHandler {
	sm := http.NewServeMux()
	db := must(sql.Open("sqlite", "./users.db"))
	ah := &AuthHandler{db, sm}
	if err := ah.initialize(ctx); err != nil {
		log.Fatalf("Cannot initialize auth: %v", err)
	}

	sm.HandleFunc("GET /auth/", func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open("auth.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}
		io.Copy(w, f)
	})
	sm.HandleFunc("GET /auth/auth.css", func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open("auth.css")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}
		io.Copy(w, f)
	})
	sm.HandleFunc("POST /auth", func(w http.ResponseWriter, r *http.Request) {
		u, pw := r.FormValue("name"), r.FormValue("password")
		rows, err := db.QueryContext(r.Context(), `SELECT id FROM users WHERE name='`+u+`' and password='`+pw+`'`)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}
		defer rows.Close()
		if !rows.Next() {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, `<html>
Invalid creadentials. <a href="/auth">Go back</a>
	</html>`)
		}
		var id int
		if err := rows.Scan(&id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		}
		http.SetCookie(w, &http.Cookie{
			Name: "userid",
			// THIS IS OF COURSE BROKEN AND A TERRIBLE AUTH MECHANISM
			// but this is a toy application and there's not benefit in
			// actually create a random token and connect that to the id
			// via DB. This application is already complicated enough as
			// is so we are taking a shortcut here.
			Value: strconv.Itoa(id),
		})
		http.Redirect(w, r, "/notes", http.StatusTemporaryRedirect)
	})
	return ah
}
