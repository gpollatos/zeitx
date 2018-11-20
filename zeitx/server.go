package zeitx

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// HTTPServer groups boilerplate needed for the http service
type HTTPServer struct {
	server *http.Server
	router *mux.Router
}

// NewHTTPServer initializes and returns a pointer to an HTTPServer based on the config
func NewHTTPServer(cfg Config) *HTTPServer {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	server := &http.Server{
		Addr:         cfg.HTTPServer.ListenAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server.Handler = handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(
		handlers.CompressHandler(Logger(JSONEnforce(router))))

	return &HTTPServer{
		server: server,
		router: router,
	}
}

// GETRoute helper
func (s *HTTPServer) GETRoute(path string, handler http.HandlerFunc) {
	s.router.HandleFunc(path, handler).Methods("GET")
}

// POSTRoute helper
func (s *HTTPServer) POSTRoute(path string, handler http.HandlerFunc) {
	s.router.HandleFunc(path, handler).Methods("POST")
}

// PATCHRoute helper
func (s *HTTPServer) PATCHRoute(path string, handler http.HandlerFunc) {
	s.router.HandleFunc(path, handler).Methods("PATCH")
}

// ListenAndServe TODO
func (s *HTTPServer) ListenAndServe() {
	log.Fatal(s.server.ListenAndServe())
}

///////////////////////////////////////////////////////////////////////////////////////////

// NotFoundHandler is a generic 404 response handler
func NotFoundHandler(w http.ResponseWriter, _ *http.Request) {
	APIError(w, http.StatusNotFound, "")
}

// MethodNotAllowedHandler is a generic 405 handler
func MethodNotAllowedHandler(w http.ResponseWriter, _ *http.Request) {
	APIError(w, http.StatusMethodNotAllowed, "")
}
