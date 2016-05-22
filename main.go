package main

import (
	"database/sql"
	"fmt"
	"html/template"
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

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("layout.html").ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		LogError(w, err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		LogError(w, err)
	}
}

func LogError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "500 Internal Server Error\n\n")
	log.Println("Server error", err)
}
