package methods

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// HTTPHandlerFunc TODO
type HTTPHandlerFunc func(interface{}, http.ResponseWriter, *http.Request)

type httpHandlerCustom struct {
	ctx interface{}
	F   func(interface{}, http.ResponseWriter, *http.Request)
}

// ServeHTTP TODO
func (h httpHandlerCustom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.F(h.ctx, w, r)
}

var diagPaths2Func = map[string]HTTPHandlerFunc{
	"/":      notImplemented,
	"/chord": chord,
}

// Server when the debug is enabled exposes a
// This data structure is protected by the Agent mutex so does not require and additional mutex here
type Server struct {
	srv               *http.Server
	mux               *http.ServeMux
	registeredHanders map[string]bool
}

// New creates a new diagnose server
func New() *Server {
	return &Server{
		registeredHanders: make(map[string]bool),
	}
}

// Init initialize the mux for the http handling and register the base hooks
func (s *Server) Init() {
	s.mux = http.NewServeMux()

	// Register local handlers
	s.RegisterHandler(s, diagPaths2Func)
}

// RegisterHandler allows to register new handlers to the mux and to a specific path
func (s *Server) RegisterHandler(ctx interface{}, hdlrs map[string]HTTPHandlerFunc) {
	for path, fun := range hdlrs {
		if _, ok := s.registeredHanders[path]; ok {
			continue
		}
		s.mux.Handle(path, httpHandlerCustom{ctx, fun})
		s.registeredHanders[path] = true
	}
}

// ServeHTTP this is the method called bu the ListenAndServe, and is needed to allow us to
// use our custom mux
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// EnableDebug opens a TCP socket to debug the passed network DB
func (s *Server) Start() {
	port := os.Getenv("PORT")
	srv := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: s}
	// Ingore ErrServerClosed that is returned on the Shutdown call
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("ListenAndServe error: %s", err)
	}
}

func notImplemented(ctx interface{}, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func ready(ctx interface{}, w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

// DebugHTTPForm helper to print the form url parameters
func DebugHTTPForm(r *http.Request) {
	for k, v := range r.Form {
		logrus.Debugf("Form[%q] = %q\n", k, v)
	}
}
