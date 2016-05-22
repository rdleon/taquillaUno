package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"LoginForm",
		"GET",
		"/login",
		LoginForm,
	},
	Route{
		"Login",
		"POST",
		"/api/login",
		Login,
	},
	Route{
		"Logout",
		"GET",
		"/api/logout",
		Logout,
	},
	Route{
		"ListEvents",
		"GET",
		"/api/events",
		ListEvents,
	},
	Route{
		"AddEvent",
		"POST",
		"/api/events",
		AddEvent,
	},
	Route{
		"UpdateEvent",
		"PUT",
		"/api/events/{eventId}",
		UpdateEvent,
	},
	Route{
		"DeleteEvent",
		"DELETE",
		"/api/events/{eventId}",
		DeleteEvent,
	},
}
