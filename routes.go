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
		"/login",
		Login,
	},
	Route{
		"Logout",
		"GET",
		"/logout",
		Logout,
	},
	Route{
		"ListEvents",
		"GET",
		"/admin/events",
		ListEvents,
	},
	Route{
		"AddEvent",
		"POST",
		"/admin/events",
		AddEvent,
	},
	Route{
		"UpdateEvent",
		"PUT",
		"/admin/events/{eventId}",
		UpdateEvent,
	},
	Route{
		"DeleteEvent",
		"DELETE",
		"/admin/events/{eventId}",
		DeleteEvent,
	},
}
