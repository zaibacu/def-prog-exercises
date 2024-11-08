package app

import "net/http"

var defaultUsers = map[string]string{
	"admin":      "admin",
	"reader":     "reader",
	"readwriter": "readwriter",
}

func Login() http.Handler {
	return nil
}
