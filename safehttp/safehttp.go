package safehttp

import (
	"log"
	"net/http"
)

type Wrapper struct {
	w http.ResponseWriter
	r *http.Request
}

func New(w http.ResponseWriter, r *http.Request) Wrapper {
	return Wrapper{w, r}
}

func (w *Wrapper) Header() http.Header {
	return w.w.Header()
}

func (w *Wrapper) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w *Wrapper) WriteHeader(statusCode int) {
	w.w.WriteHeader(statusCode)
}

func (w *Wrapper) Finalize() {

}

func (w *Wrapper) CheckSafe() bool {
	requiredHeaders := []string{"Content-Security-Policy"}

	requiredContent := make(map[string]string, 0)
	requiredContent["Content-Security-Policy"] = "script-src 'none'"

	for _, h := range requiredHeaders {
		result := w.Header().Get(h)
		if result == "" {
			log.Printf("Response if missing required header: '%s'", h)
			return false
		}

		required, ok := requiredContent[h]
		if !ok {
			// No required content
			continue
		}

		if result != required {
			log.Printf("Required header content for header: '%s': '%s', Found: '%s'", h, required, result)
			return false
		}

	}
	return true
}
