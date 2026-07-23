package rprocessor

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func reg(r *mux.Router, method, path string, handler http.Handler) {
	r.Methods(method).Path(path).Handler(handler)
}

func logWalkRoutes(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, _ := route.GetPathTemplate()
	if path == "" {
		return nil
	}
	methods, _ := route.GetMethods()
	if len(methods) == 0 {
		return nil
	}
	log.Printf("Registered route: %s %s", methods, path)
	return nil
}
