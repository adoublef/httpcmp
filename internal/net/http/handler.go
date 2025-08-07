package http

import (
	"cmp"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /ping", handlePing())
	mux.Handle("GET /users/{user}/books/{book}", handleParameters())
	mux.Handle("POST /upload", handleUpload(io.Discard))
	return mux
}

func handlePing() http.HandlerFunc {
	type response struct {
		Val string `json:"value"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, response{"pong"}, http.StatusOK)
	}
}

func handleParameters() http.HandlerFunc {
	type response struct {
		User int `json:"user"`
		Book int `json:"book"`
	}
	parse := func(_ http.ResponseWriter, r *http.Request) (user, book int, err error) {
		u, err1 := strconv.Atoi(r.PathValue("user"))
		b, err2 := strconv.Atoi(r.PathValue("book"))
		return u, b, cmp.Or(err1, err2)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user, book, err := parse(w, r)
		if err != nil {
			handleError(w, r, err, http.StatusBadRequest)
			return
		}
		respond(w, r, response{user, book}, http.StatusOK)
	}
}

func handleUpload(sink io.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mr, err := r.MultipartReader()
		if err != nil {
			handleError(w, r, err, http.StatusUnsupportedMediaType)
			return
		}
		p, err := mr.NextPart()
		if err != nil {
			handleError(w, r, err, http.StatusUnprocessableEntity)
			return
		}
		defer p.Close()
		// todo: handle filename + formname
		_, _ = io.Copy(sink, p)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error, code int) {
	// stolen from [http.Error]
	h := w.Header()
	h.Del("Content-Length")
	h.Set("X-Content-Type-Options", "nosniff")
	respond(w, r, map[string]string{"message": err.Error()}, code)
}

// respond sets application/json as Content-Type and sends the payload to client.
func respond[V any](w http.ResponseWriter, _ *http.Request, v V, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
