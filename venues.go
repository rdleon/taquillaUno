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

type Venue struct {
	VID       int    `json:"eid"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Ubication string `json:"desc"`
	Coords    string `json:"venue"`
}

// Returns the venue from the database
func getVenue(vid int) (venue Venue, err error) {
	err = db.Conn.QueryRow(
		"SELECT name, desc, ubication, coords FROM venues WHERE vid = $1",
		vid,
	).Scan(
		&(venues.Name),
		&(venues.Desc),
		&(venues.Ubication),
		&(venues.Coords),
	)

	if err != nil {
		return
	}

	venue.VID = vid
	return

}

// Creates and updates a venue in the database
func (venue Venue) Save() (err error) {
	if venue.VID > 0 {
		err = db.Conn.QueryRow(
			"INSERT INTO venues(name, desc, ubication, coords) VALUES($1, $2, $3, $4) RETURNING vid;",
			venue.Name,
			venue.Desc,
			venue.Ubication,
			venue.Coords,
		).Scan(&venue.VID)

		if err != nil {
			venues.VID = -1
		}
	} else {
		_, err = db.Conn.Query(
			"UPDATE venues SET name = $1, desc = $2, ubication = $3, coords = $4 WHERE vid = $5",
			venue.Name,
			venue.Desc,
			venue.Ubication,
			venue.Coords,
			venue.VID,
		)
	}

	return
}

// Validates all the fields in a venue
func (venue Venue) Validate() (err error) {
	return
}

// Deletes an event of the DB if it exists
func RemoveVenue(vid int) (err error) {
	_, err = db.Conn.Query("DELETE FROM venues WHERE vid = $1", vid)

	return
}

// List the first 25 venues found in the database
func listVenues() (venues []Venue, err error) {
	// TODO: Add pagination
	rows, err := db.Conn.Query(
		`SELECT vid, name, desc, ubication, coords
		 FROM venues LIMIT 25`)

	if err != nil {
		return
	}

	defer rows.Close()

	venues = make([]Venue, 0)

	for rows.Next() {
		var venue Venue
		err = rows.Scan(
			&(venue.VID),
			&(venue.Name),
			&(venue.Desc),
			&(venue.Ubication),
			&(venue.Coords),
		)

		if err != nil {
			break
		}

		venues = append(venues, venue)
	}

	return
}

func ListVenues(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	venues, err := listVenues()

	if err != nil {
		LogError(w, err)
		return
	}

	resp := map[string][]Venues{
		"venues": venues,
	}

	json.NewEncoder(w).Encode(resp)
}

func AddVenue(w http.ResponseWriter, r *http.Request) {
	var venue Venue

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&venue)

	if err != nil {
		LogError(w, err)
		return
	}

	// TODO: validate input
	err = venue.Save()

	if err != nil {
		LogError(w, err)
		return
	}

	resp := map[string]int{
		"vid": venue.VID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func UpdateVenue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if key, ok := vars["venueId"]; ok {
		var venue Venue

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&venue)

		// Get VID from request URL
		if string(venue.VID) != key {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "Bad Request"}`)
			return
		}

		if err != nil {
			LogError(w, err)
			return
		}

		// Validate input
		if err = venue.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "Bad Request"}`)
			LogError(w, err)
			return
		}

		err = venue.Save()

		if err != nil {
			LogError(w, err)
			return
		}

		resp := map[string]int{
			"vid": venue.VID,
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, `{"error": "Not Found"}`)
}

func DeleteVenue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if key, ok := vars["venueId"]; ok {
		vid, err := strconv.Atoi(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "Wrong venue id"}`)
			return
		}

		err = RemoveVenue(vid)

		if err != nil {
			LogError(w, err)
			return
		}

		resp := map[string]int{
			"deleted": vid,
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, `{"error": "Not Found"}`)
}
