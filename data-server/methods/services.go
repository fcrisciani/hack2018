package methods

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// services exposes the service to service connection data
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
			srcI, srcOK := s.serviceIPtoIndex[c.SrcIP]
			if srcOK && index != srcI {
				e.ToElement[srcI].Total++
			}
			dstI, dstOK := s.serviceIPtoIndex[c.DstIP]
			if dstOK && index != dstI {
				e.ToElement[dstI].Total++
			}
			// if both source and destination are equal to this service, this is a connection within the service
			if srcI == index && dstI == index {
				e.ToElement[index].Total++
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
