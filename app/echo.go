package app

import (
	"io"
	"net/http"
)

func Echo(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}
