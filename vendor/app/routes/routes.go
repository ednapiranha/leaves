package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"app/handlers"
)

type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
	},
	Route{
		"Profile",
		"GET",
		"/profile",
		handlers.Profile,
	},
	Route{
		"Profile",
		"POST",
		"/profile",
		handlers.Profile,
	},
	Route{
		"Directory",
		"GET",
		"/directory",
		handlers.Directory,
	},
	Route{
		"Authenticate",
		"GET",
		"/authenticate",
		handlers.Authenticate,
	},
	Route{
		"Authenticate",
		"POST",
		"/authenticate",
		handlers.Authenticate,
	},
	Route{
		"Validate",
		"GET",
		"/validate",
		handlers.Validate,
	},
	Route{
		"Validate",
		"POST",
		"/validate",
		handlers.Validate,
	},
	Route{
		"Logout",
		"GET",
		"/logout",
		handlers.Logout,
	},
}


func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}