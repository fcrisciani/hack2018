package methods

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fcrisciani/hack2018/data-server/elastic"
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
	"/":         notImplemented,
	"/chord":    chord,
	"/services": services,
	"/pods":     pods,
}

// Server when the debug is enabled exposes a
// This data structure is protected by the Agent mutex so does not require and additional mutex here
type Server struct {
	srv               *http.Server
	mux               *http.ServeMux
	registeredHanders map[string]bool

	serviceList      []*ServiceConnections
	serviceIPtoIndex map[string]int
	serviceListLock  sync.Mutex
	podList          []*PodConnections
	podIPtoIndex     map[string]int
	podListLock      sync.Mutex
}

// New creates a new diagnose server
func New() *Server {
	return &Server{
		registeredHanders: make(map[string]bool),
		serviceIPtoIndex:  make(map[string]int),
		podIPtoIndex:      make(map[string]int),
	}
}

// Init initialize the mux for the http handling and register the base hooks
func (s *Server) Init() {
	// initialize services data
	s.refreshServiceList()
	s.refreshPodList()

	// keeps data fresh
	// go s.refreshData()

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

// Start starts the server
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

type ServiceConnections struct {
	service     *elastic.Service
	pods        []*elastic.Pod
	connections []*elastic.Connection
}

type PodConnections struct {
	pod         *elastic.Pod
	connections []*elastic.Connection
}

func (s *Server) refreshServiceList() {
	log := logrus.WithField("method", "refreshServices")
	log.Debug("start")

	services, err := elastic.GetServices()
	if err != nil {
		log.WithError(err).Error("error in getting services")
		panic(err)
	}

	s.serviceListLock.Lock()
	newServices := make(map[string]bool)
	for _, s := range services {
		newServices[s.Meta.ServiceName] = true
	}

	// Check if service list changed
	for _, srv := range s.serviceList {
		if !newServices[srv.service.Meta.ServiceName] {
			log.Infof("service %s is gone", srv.service.Meta.ServiceName)
			// TODO delete it
		}
	}

	var foundNew bool
	for _, srv := range services {
		if _, ok := s.serviceIPtoIndex[srv.Spec.ClusterIP]; !ok {
			log.Infof("service %s is new", srv.Meta.ServiceName)
			s.serviceList = append(s.serviceList, &ServiceConnections{srv, nil, nil})
			s.serviceIPtoIndex[srv.Spec.ClusterIP] = len(s.serviceList) - 1
			foundNew = true
		}
	}

	if foundNew {
		s.refreshServicePodList()
		s.refreshServiceConnections()
	}

	s.serviceListLock.Unlock()

	log.Debug("done")
}

func (s *Server) refreshServicePodList() {
	log := logrus.WithField("method", "refreshPods")
	log.Debug("start")

	for index, srv := range s.serviceList {
		pods, err := elastic.GetPodsForService(srv.service)
		if err != nil {
			log.WithError(err).Error("error in getting pods for service")
			continue
		}
		srv.pods = pods
		// all the pod IP are part of the service update the list
		for _, p := range pods {
			s.serviceIPtoIndex[p.Status.PodIP] = index
			log.Infof("Pod: %v IP:%v part of %v", p.Meta.Name, p.Status.PodIP, srv.service.Meta.ServiceName)
		}
	}
	log.Debug("done")
}

func (s *Server) refreshServiceConnections() {
	log := logrus.WithField("method", "refreshConnections")
	log.Debug("start")
	for _, srv := range s.serviceList {
		connections, err := elastic.GetAllConnections(srv.service.Spec.ClusterIP, 0)
		if err != nil {
			log.WithError(err).Error("error in getting connections")
			panic(err)
		}
		srv.connections = connections
		for _, pod := range srv.pods {
			connections, err := elastic.GetAllConnections(pod.Status.PodIP, 0)
			if err != nil {
				log.WithError(err).Error("error in getting connections")
				continue
			}
			srv.connections = append(srv.connections, connections...)
		}
		log.Infof("service %s has %v connections", srv.service.Meta.ServiceName, len(srv.connections))
	}
	log.Debug("done")
}

func (s *Server) refreshPodList() {
	log := logrus.WithField("method", "refreshPods")
	log.Debug("start")

	pods, err := elastic.GetPods()
	if err != nil {
		log.WithError(err).Error("error in getting pods")
		panic(err)
	}

	newPods := make(map[string]bool)
	for _, p := range pods {
		newPods[p.Meta.Name] = true
	}

	s.podListLock.Lock()
	// Check if service list changed
	for _, p := range s.podList {
		if !newPods[p.pod.Meta.Name] {
			log.Infof("pod %s is gone", p.pod.Meta.Name)
			// TODO delete it
		}
	}

	var foundNew bool
	for _, pod := range pods {
		if _, ok := s.podIPtoIndex[pod.Status.PodIP]; !ok {
			log.Infof("pod %s is new", pod.Meta.Name)
			s.podList = append(s.podList, &PodConnections{pod, nil})
			s.podIPtoIndex[pod.Status.PodIP] = len(s.podList) - 1
			foundNew = true
		}
	}

	if foundNew {
		s.refreshPodConnections()
	}

	s.podListLock.Unlock()
	log.Debug("done")
}

func (s *Server) refreshPodConnections() {
	log := logrus.WithField("method", "refreshPodConnections")
	log.Debug("start")
	for _, pod := range s.podList {
		connections, err := elastic.GetAllConnections(pod.pod.Status.PodIP, 0)
		if err != nil {
			log.WithError(err).Error("error in getting connection for pod")
			continue
		}
		pod.connections = connections
	}
	log.Debug("done")
}

// refreshServices runs periodically and keeps the service matrix fresh
func (s *Server) refreshData() {
	log := logrus.WithField("method", "refreshServices")
	log.Debug("started")
	iteration := 1
	const FULLREFRESH int = 10
	for {
		// sleep for a random time between 0 and 5s
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		log.Info("new cycle")
		if iteration == 0 {
			s.refreshServiceList()
		} else if iteration == 5 {
			s.refreshPodList()
		} else {
			if iteration%2 == 0 {
				s.serviceListLock.Lock()
				s.refreshServicePodList()
				s.refreshServiceConnections()
				s.serviceListLock.Unlock()
			} else {
				s.podListLock.Lock()
				s.refreshPodConnections()
				s.podListLock.Unlock()
			}
		}

		iteration = (iteration + 1) % FULLREFRESH
		log.Debug("done")
	}
}
