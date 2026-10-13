package rprocessor

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
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
	log.Info().Strs("methods", methods).Str("path", path).Msg("Registered route")
	return nil
}
