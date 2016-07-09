package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

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
		"Profile",
		"GET",
		"/profile",
		Profile,
	},
	Route{
		"Profile",
		"POST",
		"/profile",
		Profile,
	},
	Route{
		"Directory",
		"GET",
		"/directory",
		Directory,
	},
	Route{
		"Directory",
		"GET",
		"/search",
		Directory,
	},
	Route{
		"Search",
		"POST",
		"/search",
		Directory,
	},
	Route{
		"Strain",
		"GET",
		"/strain/{ucpc}",
		StrainDetail,
	},
	Route{
		"UpdateStrain",
		"GET",
		"/add/strain",
		UpdateStrain,
	},
	Route{
		"UpdateStrain",
		"POST",
		"/add/strain",
		UpdateStrain,
	},
	Route{
		"Review",
		"POST",
		"/strain/{ucpc}",
		UpdateReview,
	},
	Route{
		"GetReview",
		"GET",
		"/review/{rid}",
		GetReview,
	},
	Route{
		"DeleteReview",
		"GET",
		"/delete/{rid}",
		DeleteReview,
	},
	Route{
		"UpdateLike",
		"GET",
		"/like/{rid}",
		UpdateLike,
	},
	Route{
		"Authenticate",
		"GET",
		"/authenticate",
		Authenticate,
	},
	Route{
		"Authenticate",
		"POST",
		"/authenticate",
		Authenticate,
	},
	Route{
		"Validate",
		"GET",
		"/validate",
		Validate,
	},
	Route{
		"Validate",
		"POST",
		"/validate",
		Validate,
	},
	Route{
		"Logout",
		"GET",
		"/logout",
		Logout,
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
