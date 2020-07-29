package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kerti/evm/02-kitara-store/handler/response"
	"github.com/kerti/evm/02-kitara-store/util/logger"
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

	// Orders
	s.router.HandleFunc("/orders/process", s.OrderHandler.HandleProcessOrder).Methods("POST")

	http.Handle("/", s.router)
}
