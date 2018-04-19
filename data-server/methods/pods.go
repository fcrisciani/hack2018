package methods

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// pods exposes the pod to pod connection data
func pods(ctx interface{}, w http.ResponseWriter, r *http.Request) {
	log := logrus.WithField("method", "pods")
	log.Info("new request")

	s := ctx.(*Server)
	c := &chordGraph{
		Data: make([]Element, 0, len(s.serviceList)),
	}

	s.podListLock.Lock()

	for ip, index := range s.podIPtoIndex {
		log.Infof("%v --> %d --> %s", ip, index, s.podList[index].pod.Meta.Name)
	}

	for index, p := range s.podList {
		e := Element{Name: p.pod.Meta.Name, IP: p.pod.Status.PodIP, ToElement: make([]Connection, len(s.podList))}
		// check connections
		for _, c := range p.connections {
			log.Infof("processing connection %+v from %v", c, p.pod.Meta.Name)
			srcI, srcOK := s.podIPtoIndex[c.SrcIP]
			if srcOK && index != srcI {
				e.ToElement[srcI].Total++
			}
			dstI, dstOK := s.podIPtoIndex[c.DstIP]
			if dstOK && index != dstI {
				e.ToElement[dstI].Total++
			}
			// if both source and destination are equal to this service, this is a connection within the pod
			if srcI == index && dstI == index {
				e.ToElement[index].Total++
			}
		}
		log.Infof("pod %s(%s) row:%+v", p.pod.Meta.Name, p.pod.Status.PodIP, e)
		// append the row
		c.Data = append(c.Data, e)
	}

	s.podListLock.Unlock()

	log.Infof("final body:%+v", c)

	res, _ := json.Marshal(c)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "%s", res)
}
