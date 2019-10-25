package utils

import (
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var m = sync.Mutex{}

func startHTTPServer(callback func(r http.Request)) *http.Server {
	r := mux.NewRouter()
	r.Handle("/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		callback(*r)
		m.Unlock()
		io.WriteString(w, "hello world\n")
	}))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go srv.ListenAndServe()

	// returning reference so caller can call Shutdown()
	return srv
}
