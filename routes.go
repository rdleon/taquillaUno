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
		"/",
		false,
		Index,
	},
	Route{
		"Login",
		"POST",
		"/login",
		false,
		Login,
	},
	Route{
		"Logout",
		"GET",
		"/logout",
		false,
		Logout,
	},
	Route{
		"CreateUser",
		"POST",
		"/users",
		true,
		CreateUser,
	},
	Route{
		"ListEvents",
		"GET",
		"/events",
		true,
		ListEvents,
	},
	Route{
		"AddEvent",
		"POST",
		"/events",
		true,
		AddEvent,
	},
	Route{
		"UpdateEvent",
		"PUT",
		"/events/{eventId}",
		true,
		UpdateEvent,
	},
	Route{
		"DeleteEvent",
		"DELETE",
		"/events/{eventId}",
		true,
		DeleteEvent,
	},
}
