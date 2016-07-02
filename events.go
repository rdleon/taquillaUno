package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rdleon/taquillaUno/db"
)

type Event struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Date      time.Time `json:"date"`
	Duration  int       `json:"duraction"`
	Published bool      `json:"published"`
}

type Events []Event

func listEvents(page int) Events {
	rows, err := db.Conn.Query("SELECT name, date, duration, published FROM events ORDER BY published desc LIMIT 25")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var events Events = make([]Event, 0)

	for rows.Next() {
		var event Event
		err = rows.Scan(&(event.Name), &(event.Date), &(event.Duration), &(event.Published))
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
