package api

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	auth2 "shopingList/api/auth"
	"shopingList/store"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// ServerAPI - interface for REST API implementation
type ServerAPI interface {
	Routes() []Route
}

// Route object for rest server
type Route struct {
	Name   string
	Method string
	Path   string
	Func   http.HandlerFunc
}

// Rest server
type Rest struct {
	authenticator *auth2.Service
	dataService   store.DataService
	port          string
	server        *http.Server
	stop          chan struct{}
	router        mux.Router
}

// New returns new instance of REST server
func New(a *auth2.Service, d store.DataService, port string, stop chan struct{}) *Rest {
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter() // TODO later create versioning support
	return &Rest{authenticator: a, dataService: d, port: port, stop: stop, router: *router}
}

// Run rest server
func (s *Rest) Run() {
	s.server = &http.Server{
		Addr:         ":" + s.port,
		Handler:      &s.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info("starting server at ", s.server.Addr)

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("HTTP server ListenAndServe error: %v", err)
		close(s.stop)
	}
}

// Stop rest server
func (s *Rest) Stop() {
	log.Info("stopping server")

	if s.server == nil {
		log.Error("server is nil")
		return
	}
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Error("HTTP server Shutdown error: %v", err)
	}
}

func (s *Rest) AddPublicRoutes(routes ...Route) {
	for _, route := range routes {
		s.router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.Func)
	}
}

func (s *Rest) AddPrivateRoutes(routes ...Route) {
	for _, route := range routes {
		handler := alice.New(s.accessTokenValidationMiddleware).ThenFunc(route.Func)
		s.router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(handler)
	}
}
