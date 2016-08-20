package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/rdleon/taquillaUno/db"

	_ "github.com/lib/pq"
)

var Conf map[string]string

func main() {
	var err error

	// TODO: Load Configuration from file
	Conf = map[string]string{
		"db_user":   "taquilla",
		"db_passwd": "secret",
		"db_name":   "taquilla",
		"db_host":   "localhost",
		"listen":    "127.0.0.1",
		"secret":    "secretkeyVerySecret",
	}

	dbConf := fmt.Sprintf("user=%s dbname=%s password=%s", Conf["db_user"], Conf["db_name"], Conf["db_passwd"])
	db.Conn, err = sql.Open("postgres", dbConf)

	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}

	defer db.Conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	// TODO: Set public route
	router := NewRouter()

	fmt.Println("Runnig server on " + Conf["listen"])
	log.Fatal(http.ListenAndServe(Conf["listen"], router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"version\": 1}")
}

func LogError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "{\"error\": \"Internal Server Error\"}")
	log.Println("Server error", err)
}
