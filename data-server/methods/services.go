package methods

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/fcrisciani/hack2018/data-server/elastic"
	"github.com/sirupsen/logrus"
)

type ServiceConnections struct {
	service     *elastic.Service
	pods        []*elastic.Pod
	connections []*elastic.Connection
}

// initServices build the first matrix, no need to lock because there is no contention
// happens before the registration of the endpoint
func (s *Server) initServices() {
	log := logrus.WithField("method", "initServices")
	log.Info("initialize services")

	s.refreshServices()
	s.refreshPods()
	s.refreshConnections()

	go s.refreshData()
	log.Info("done")
}

func (s *Server) refreshServices() {
	log := logrus.WithField("method", "refreshServices")
	log.Info("start")

	services, err := elastic.GetServices()
	if err != nil {
		log.WithError(err).Error("error in getting services")
		panic(err)
	}

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

	for _, srv := range services {
		if _, ok := s.serviceIPtoIndex[srv.Spec.ClusterIP]; !ok {
			log.Infof("service %s is new", srv.Meta.ServiceName)
			s.serviceList = append(s.serviceList, &ServiceConnections{srv, nil, nil})
			s.serviceIPtoIndex[srv.Spec.ClusterIP] = len(s.serviceList) - 1
		}
	}

	log.Info("done")
}

func (s *Server) refreshPods() {
	log := logrus.WithField("method", "refreshPods")
	log.Info("start")
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
	log.Info("done")
}

func (s *Server) refreshConnections() {
	log := logrus.WithField("method", "refreshConnections")
	log.Info("start")
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
	log.Info("done")
}

// refreshServices runs periodically and keeps the service matrix fresh
func (s *Server) refreshData() {
	log := logrus.WithField("method", "refreshServices")
	log.Info("started")
	iteration := 1
	const FULLREFRESH int = 10
	for {
		// sleep for a random time between 0 and 5s
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		log.Info("new cycle")
		s.serviceListLock.Lock()
		if iteration == 0 {
			s.refreshServices()
		}
		s.refreshPods()
		s.refreshConnections()

		s.serviceListLock.Unlock()

		iteration = (iteration + 1) % FULLREFRESH
		log.Info("done")
	}
}

func services(ctx interface{}, w http.ResponseWriter, r *http.Request) {
	log := logrus.WithField("method", "services")
	log.Info("new request")

	s := ctx.(*Server)
	c := &chordGraph{
		Data: make([]Element, 0, len(s.serviceList)),
	}

	s.serviceListLock.Lock()

	for ip, index := range s.serviceIPtoIndex {
		log.Infof("%v --> %d --> %s", ip, index, s.serviceList[index].service.Meta.ServiceName)
	}

	for index, srv := range s.serviceList {
		e := Element{Name: srv.service.Meta.ServiceName, IP: srv.service.Spec.ClusterIP, ToElement: make([]Connection, len(s.serviceList))}
		// check connections
		for _, c := range srv.connections {
			log.Infof("processing connection %+v from %v", c, srv.service.Meta.ServiceName)
			// log.Infof("%s index:%d", c.SrcIP, s.serviceIPtoIndex[c.DstIP])
			if i, ok := s.serviceIPtoIndex[c.SrcIP]; ok && index != i {
				// log.Infof("incrementing %d", i)
				e.ToElement[i].Total++
			}
			if i, ok := s.serviceIPtoIndex[c.DstIP]; ok && index != i {
				// log.Infof("incrementing %d", i)
				e.ToElement[i].Total++
			}
		}
		log.Infof("service %s(%s) row:%+v", srv.service.Meta.ServiceName, srv.service.Spec.ClusterIP, e)
		// append the row
		c.Data = append(c.Data, e)
	}

	s.serviceListLock.Unlock()

	log.Infof("final body:%+v", c)

	res, _ := json.Marshal(c)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "%s", res)
}
