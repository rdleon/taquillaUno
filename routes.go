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
		"LoginForm",
		"GET",
		"/login",
		false,
		LoginForm,
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
		true,
		Logout,
	},
	Route{
		"ListEvents",
		"GET",
		"/admin/events",
		true,
		ListEvents,
	},
	Route{
		"AddEvent",
		"POST",
		"/admin/events",
		true,
		AddEvent,
	},
	Route{
		"UpdateEvent",
		"PUT",
		"/admin/events/{eventId}",
		true,
		UpdateEvent,
	},
	Route{
		"DeleteEvent",
		"DELETE",
		"/admin/events/{eventId}",
		true,
		DeleteEvent,
	},
}
