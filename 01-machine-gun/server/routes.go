package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kerti/evm/01-machine-gun/handler/response"
	"github.com/kerti/evm/01-machine-gun/util/logger"
)

// InitRoutes initializes the routes
func (s *Server) InitRoutes() {
	logger.Trace("Initializing routes...")
	s.router = mux.NewRouter()

	// Preflight for CORS
	s.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.RespondWithNoContent(w)
	})

	// Health
	s.router.HandleFunc("/health", s.HealthHandler.HandleHealthCheck).Methods("GET")

	http.Handle("/", s.router)
}
