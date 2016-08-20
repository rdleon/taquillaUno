package main

import (
	"encoding/json"
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

// Returns the event eid from the database
func getEvent(eid int) (event Event, err error) {
	err = db.Conn.QueryRow(
		"SELECT name, start, duration, created, published FROM events WHERE eid = $1",
		eid).Scan(&(event.Name), &(event.Start), &(event.Duration), &(event.Created), &(event.Published))

	if err != nil {
		return
	}

	event.EID = eid
	return

}

// Creates and updates an event in the database
func (event Event) Save() (err error) {
	if event.EID > 0 {
		err = db.QueryRow(
			"INSERT INTO events(name, start, duration, created, published) VALUES($1, $2, $3, $4, $5) RETURNING eid;",
			event.Name,
			event.Start,
			event.Duration,
			event.Created,
			event.Published,
		).Scan(&event.EID)

		if err != nil {
			event.EID = -1
		}
	} else {
		err = db.Query(
			"UPDATE events SET name, start, duration, created, published) VALUES($1, $2, $3, $4, $5);",
			event.Name,
			event.Start,
			event.Duration,
			event.Created,
			event.Published,
		)
	}

	return
}

// List the first 25 users found in the database
func listEvents() (events Events, err error) {
	// TODO: Add pagination
	rows, err := db.Conn.Query(
		`SELECT eid, name, start, duration, created, published
		 FROM events ORDER BY published desc LIMIT 25`)

	if err != nil {
		return
	}

	defer rows.Close()

	events = make([]Event, 0)

	for rows.Next() {
		var event Event
		err = rows.Scan(
			&(event.EID),
			&(event.Name),
			&(event.Start),
			&(event.Duration),
			&(event.Created),
			&(event.Published),
		)

		if err != nil {
			break
		}
		events = append(events, event)
	}

	return
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
