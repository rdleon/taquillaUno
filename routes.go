package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	AuthNeeded  bool
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api",
		false,
		Index,
	},
	Route{
		"Login",
		"POST",
		"/api/login",
		false,
		Login,
	},
	Route{
		"Logout",
		"GET",
		"/api/logout",
		false,
		Logout,
	},
	Route{
		"CreateUser",
		"POST",
		"/api/users",
		true,
		CreateUser,
	},
	Route{
		"ListEvents",
		"GET",
		"/api/events",
		true,
		ListEvents,
	},
	Route{
		"AddEvent",
		"POST",
		"/api/events",
		true,
		AddEvent,
	},
	Route{
		"UpdateEvent",
		"PUT",
		"/api/events/{eventId}",
		true,
		UpdateEvent,
	},
	Route{
		"DeleteEvent",
		"DELETE",
		"/api/events/{eventId}",
		true,
		DeleteEvent,
	},
	Route{
		"ListVenues",
		"GET",
		"/api/venues",
		true,
		ListVenues,
	},
	Route{
		"AddVenues",
		"POST",
		"/api/venues",
		true,
		AddVenue,
	},
	Route{
		"UpdateVenue",
		"PUT",
		"/api/venues/{venueId}",
		true,
		UpdateVenue,
	},
	Route{
		"DeleteVenue",
		"DELETE",
		"/api/venues/{venueId}",
		true,
		DeleteVenue,
	},
}
