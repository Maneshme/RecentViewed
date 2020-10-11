package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//Route ...
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func HandleRequests() {
	router := BuildRouter(BuildRoutes())
	router.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swaggerui/"))))
	http.Handle("/", router)
	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		log.Fatal(err)
	}
	// log.Fatal(http.ListenAndServe(":10000", nil))
}

func BuildRoutes() Routes {
	var routes = Routes{
		Route{
			Name:        "HealthCheck",
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: homePage,
		},
		Route{
			Name:        "Recent",
			Method:      "GET",
			Pattern:     "/recent",
			HandlerFunc: getRecentViewed,
		},
		Route{
			Name:        "Recent",
			Method:      "POST",
			Pattern:     "/recent",
			HandlerFunc: postRecentViewed,
		},
	}
	return routes
}

//BuildRouter Builds a Mux router from the given route definitions
func BuildRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// router.Use(loggingMiddleware(router))

	for _, route := range routes {
		router.
			Methods(strings.Split(route.Method, ",")...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
