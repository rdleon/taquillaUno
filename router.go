package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func checkForAuth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			sess *sessions.Session
			err  error
		)

		sess, err = Store.Get(r, "logged")

		if err != nil {
			LogError(w, err)
			return
		}

		tmp := sess.Values["uid"]

		if tmp != nil {
			// Already logged in
			inner.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401 Unauthorized\n\n")
		}
	})
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		if route.AuthNeeded {
			handler = checkForAuth(handler)
		}
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Serve the statics assets of the sites (imgs, css and js)
	router.
		Methods("GET").
		PathPrefix("/static/").
		Name("Public files").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	return router
}
