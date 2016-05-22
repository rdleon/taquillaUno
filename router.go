package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
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
