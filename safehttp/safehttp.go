package safehttp

import "net/http"

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
	return true
}
