package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kerti/evm/04-tennis-player/handler/response"
	"github.com/kerti/evm/04-tennis-player/util/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

// InitRoutes initializes the routes
func (s *Server) InitRoutes() {
	logger.Trace("Initializing routes...")
	s.router = mux.NewRouter()

	// Preflight for CORS
	s.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.RespondWithNoContent(w)
	})

	// Swagger Docs
	s.router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	// Health
	s.router.HandleFunc("/health", s.HealthHandler.HandleHealthCheck).Methods("GET")

	// Player
	s.router.HandleFunc("/players", s.PlayerHandler.HandleCreate).Methods("POST")
	s.router.HandleFunc("/players/addBall", s.PlayerHandler.HandleAddBall).Methods("POST")

	http.Handle("/", s.router)
}
