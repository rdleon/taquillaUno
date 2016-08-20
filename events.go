package main

import (
	"encoding/json"
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

func listEvents(page int) (events Events, err error) {
	rows, err := db.Conn.Query("SELECT eid, name, start, duration, created, published" +
		" FROM events ORDER BY published desc LIMIT 25",
	)

	if err != nil {
		return
	}

	defer rows.Close()

	events = make([]Event, 0)

	for rows.Next() {
		var event Event
		err = rows.Scan(&(event.EID), &(event.Name), &(event.Start), &(event.Duration), &(event.Created), &(event.Published))
		if err != nil {
			break
		}
		events = append(events, event)
	}

	return
}

func loadEvent(id int) {
}

func saveEvent(event Event) {
}

func ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := listEvents()

	if err != nil {
		LogError(w, err)
		return
	}

	resp := map[string]Events{
		"events": events,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(resp)
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
}
