package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rdleon/taquillaUno/db"
)

type Event struct {
	EID       int       `json:"eid"`
	Name      string    `json:"name"`
	Start     time.Time `json:"start"`
	Duration  int       `json:"duration"`
	Created   time.Time `json:"created"`
	Published bool      `json:"published"`
}

type Events []Event

func listEvents(page int) Events {
	rows, err := db.Conn.Query("SELECT name, start, duration, published" +
		" FROM events ORDER BY published desc LIMIT 25",
	)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var events Events = make([]Event, 0)

	for rows.Next() {
		var event Event
		err = rows.Scan(&(event.Name), &(event.Start), &(event.Duration), &(event.Published))
		events = append(events, event)
	}

	return events
}

func loadEvent(id int) {
}

func saveEvent(event Event) {
}

func ListEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
}
