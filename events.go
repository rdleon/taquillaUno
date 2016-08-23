package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rdleon/taquillaUno/db"
)

type Event struct {
	EID       int       `json:"eid"`
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	Venue     string    `json:"venue"`
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
		eid,
	).Scan(
		&(event.Name),
		&(event.Start),
		&(event.Duration),
		&(event.Created),
		&(event.Published),
	)

	if err != nil {
		return
	}

	event.EID = eid
	return

}

// Creates and updates an event in the database
func (event Event) Save() (err error) {
	if event.EID > 0 {
		err = db.Conn.QueryRow(
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
		_, err = db.Conn.Query(
			"UPDATE events SET name = $1, start = $2, duration = $3, created = $4, published = $5 WHERE eid = $6",
			event.Name,
			event.Start,
			event.Duration,
			event.Created,
			event.Published,
			event.EID,
		)
	}

	return
}

// Validates all the fields in an event
func (event Event) Validate() (err error) {
	return
}

// Deletes an event of the DB if it exists
func RemoveEvent(eid int) (err error) {
	_, err = db.Conn.Query("DELETE FROM events WHERE eid = $1", eid)

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
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	events, err := listEvents()

	if err != nil {
		LogError(w, err)
		return
	}

	resp := map[string]Events{
		"events": events,
	}

	json.NewEncoder(w).Encode(resp)
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	var event Event

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&event)

	if err != nil {
		LogError(w, err)
		return
	}

	// TODO: validate input
	err = event.Save()

	if err != nil {
		LogError(w, err)
		return
	}

	resp := map[string]int{
		"eid": event.EID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// TODO: Add get (single) event

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if key, ok := vars["eventId"]; ok {
		var event Event

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&event)

		// Get EID from request URL
		if string(event.EID) != key {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "Bad Request"}`)
			return
		}

		if err != nil {
			LogError(w, err)
			return
		}

		// Validate input
		if err = event.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "Bad Request"}`)
			LogError(w, err)
			return
		}

		err = event.Save()

		if err != nil {
			LogError(w, err)
			return
		}

		resp := map[string]int{
			"eid": event.EID,
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, `{"error": "Not Found"}`)
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if key, ok := vars["eventId"]; ok {
		eid, err := strconv.Atoi(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "Wrong event id"}`)
			return
		}

		err = RemoveEvent(eid)

		if err != nil {
			LogError(w, err)
			return
		}

		resp := map[string]int{
			"deleted": eid,
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, `{"error": "Not Found"}`)
}
