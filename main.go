package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rdleon/taquillaUno/db"

	_ "github.com/lib/pq"
)

var Store *sessions.CookieStore

func main() {
	var err error

	Store = sessions.NewCookieStore([]byte("secretsecret"))

	db.Conn, err = sql.Open("postgres", "user=taquilla dbname=taquilla password=secret")

	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}

	defer db.Conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	// TODO: Set public route
	router := NewRouter()

	fmt.Println("Runnig server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"version\": 1}")
}

func LogError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "{\"status\": \"error\", \"message\": \"Internal Server Error\"}")
	log.Println("Server error", err)
}
