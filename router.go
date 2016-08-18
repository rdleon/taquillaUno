package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func CheckAuth(r *http.Request) (string, bool) {
	tokenStr := r.Header.Get("Authorization")

	if tokenStr == "" {
		return "", false
	}

	i := strings.Index(tokenStr, "Bearer ")

	if i >= 0 {
		tokenStr = tokenStr[i+len("Bearer "):]
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}

		// TODO: change this key!!!
		return []byte("verysecretKey"), nil
	})

	if err != nil {
		return "", false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), true
	}

	return "", false
}

func checkForAuth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if _, ok := CheckAuth(r); ok {
			// TODO: include parsed user info in request??
			inner.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "401 Unauthorized\n\n")
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

	return router
}
