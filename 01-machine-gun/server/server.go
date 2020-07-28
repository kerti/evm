package server

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/kerti/evm/01-machine-gun/config"
	"github.com/kerti/evm/01-machine-gun/handler"
	"github.com/kerti/evm/01-machine-gun/util/logger"

	"github.com/gorilla/mux"
)

var isShuttingDown bool

// Server is the server instance
type Server struct {
	config        *config.Config
	HealthHandler handler.Health `inject:"healthHandler"`
	router        *mux.Router
}

// Startup perform startup functions
func (s *Server) Startup() {
	logger.Trace("HTTP Server starting up...")
	s.config = config.Get()
	s.InitRoutes()
	s.InitMiddleware()
}

// Shutdown cleans up everything and shuts down
func (s *Server) Shutdown() {
	logger.Trace("HTTP Server shutting down...")
}

// InitMiddleware initializes all middlewares
func (s *Server) InitMiddleware() {
	s.router.Use(s.loggingMiddleware)
}

// Serve runs the server
func (s *Server) Serve() {
	logger.Info("Server started and is available at the following address(es):")
	ips, _ := s.getIPs()
	for _, ip := range ips {
		logger.Info("- http://%s:%d", ip.String(), s.config.Server.Port)
	}

	logger.Fatal("%s", http.ListenAndServe(fmt.Sprintf(":%d", s.config.Server.Port), nil))

}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerList := make([]string, 0)
		for k, v := range r.Header {
			headerList = append(headerList, fmt.Sprintf("  - %s: %s", k, v))
		}
		headers := fmt.Sprintf("- HEADERS:\n%s", strings.Join(headerList, "\n"))

		cookieList := make([]string, 0)
		for _, v := range r.Cookies() {
			cookieList = append(cookieList, fmt.Sprintf(" - %s: %v", v.Name, v.Value))
		}
		cookies := fmt.Sprintf("- COOKIES:\n%s", strings.Join(cookieList, "\n"))

		logger.Trace("### RECEIEVED %v %v\n%s\n%s", r.Method, r.RequestURI, headers, cookies)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) getIPs() ([]net.IP, error) {
	res := make([]net.IP, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						res = append(res, ipNet.IP)
					}
				}
			}
		}
	}
	return res, nil
}
