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
				e.ToElement[srcI].Flows = append(e.ToElement[srcI].Flows, *c)
			}
			dstI, dstOK := s.serviceIPtoIndex[c.DstIP]
			if dstOK && index != dstI {
				e.ToElement[dstI].Total++
				e.ToElement[dstI].Flows = append(e.ToElement[dstI].Flows, *c)
			}
			// if both source and destination are equal to this service, this is a connection within the service
			if srcOK && dstOK && srcI == index && dstI == index {
				e.ToElement[index].Total++
				e.ToElement[index].Flows = append(e.ToElement[index].Flows, *c)
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

// servicesToUnknown exposes the service to service connection data
func servicesToUnknown(ctx interface{}, w http.ResponseWriter, r *http.Request) {
	log := logrus.WithField("method", "services")
	log.Info("new request")

	s := ctx.(*Server)
	c := &chordGraph{
		Data: make([]Element, 0, len(s.serviceList)),
	}

	s.serviceListLock.Lock()
	var unknownIP []string
	unknownIPtoIndex := make(map[string]int)
	for _, srv := range s.serviceList {
		e := Element{Name: srv.service.Meta.ServiceName, IP: srv.service.Spec.ClusterIP, ToElement: nil}
		for _, c := range srv.connections {
			if _, ok := s.serviceIPtoIndex[c.SrcIP]; !ok {
				if _, uOk := unknownIPtoIndex[c.SrcIP]; !uOk {
					unknownIP = append(unknownIP, c.SrcIP)
					unknownIPtoIndex[c.SrcIP] = len(unknownIP) - 1
				}
			}
			if _, ok := s.serviceIPtoIndex[c.DstIP]; !ok {
				if _, uOk := unknownIPtoIndex[c.DstIP]; !uOk {
					unknownIP = append(unknownIP, c.DstIP)
					unknownIPtoIndex[c.DstIP] = len(unknownIP) - 1
				}
			}
		}
		c.Data = append(c.Data, e)
	}
	log.Infof("unknown IPs: %v", unknownIP)

	for _, ip := range unknownIP {
		c.Data = append(c.Data, Element{Name: ip, IP: ip, ToElement: make([]Connection, len(s.serviceList)+len(unknownIPtoIndex))})
	}

	// FIXME
	for i, data := range c.Data {
		log.Infof("%d) %v", i, data)
		if len(data.ToElement) == 0 {
			data.ToElement = make([]Connection, len(s.serviceList)+len(unknownIPtoIndex))
		}
		c.Data[i] = data
	}

	// mark connections for services
	for index, srv := range s.serviceList {
		e := c.Data[index]

		for _, con := range srv.connections {
			log.Infof("processing connection %+v from %v", con, srv.service.Meta.ServiceName)
			// log.Infof("%s index:%d", c.SrcIP, s.serviceIPtoIndex[c.DstIP])
			srcI, srcOK := unknownIPtoIndex[con.SrcIP]
			if srcOK {
				log.Infof("Increment for %v", e.Name)
				e.ToElement[len(s.serviceList)+srcI].Total++
				e.ToElement[len(s.serviceList)+srcI].Flows = append(e.ToElement[len(s.serviceList)+srcI].Flows, *con)
				log.Infof("Increment for %v", c.Data[len(s.serviceList)+srcI].Name)
				c.Data[len(s.serviceList)+srcI].ToElement[index].Total++
			}
			dstI, dstOK := unknownIPtoIndex[con.DstIP]
			if dstOK {
				log.Infof("Increment for %v", e.Name)
				// log.Infof("index:%d dstI:%d element:%d and Data:%d", len(s.serviceList)+dstI-1, dstI, len(e.ToElement), len(c.Data))
				e.ToElement[len(s.serviceList)+dstI].Total++
				log.Infof("Increment for %v", c.Data[len(s.serviceList)+dstI].Name)
				e.ToElement[len(s.serviceList)+dstI].Flows = append(e.ToElement[len(s.serviceList)+dstI].Flows, *con)
				// log.Infof("%+v", c.Data[len(s.serviceList)+dstI-1])
				// log.Infof("%+v", c.Data[len(s.serviceList)+dstI-1].ToElement[dstI])
				c.Data[len(s.serviceList)+dstI].ToElement[index].Total++
			}
		}
		log.Infof("service %s(%s) row:%+v", srv.service.Meta.ServiceName, srv.service.Spec.ClusterIP, e)
	}

	s.serviceListLock.Unlock()
	//
	// log.Infof("final body:%+v", c)

	res, _ := json.Marshal(c)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "%s", res)
}
