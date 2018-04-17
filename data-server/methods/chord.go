package methods

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Connection struct {
	Total int32 `json:"total"`
}

type Element struct {
	Name      string       `json:"name"`
	IP        string       `json:"ip"`
	ToElement []Connection `json:"connections"`
}

type chordGraph struct {
	Data []Element `json:"graph"`
}

func chord(ctx interface{}, w http.ResponseWriter, r *http.Request) {
	logrus.WithField("method", "chord").Info("new request")
	c := &chordGraph{
		Data: []Element{
			{Name: "service 1", IP: "192.168.0.1", ToElement: []Connection{{0}, {rand.Int31n(10)}, {rand.Int31n(10)}, {rand.Int31n(10)}}},
			{Name: "service 2", IP: "192.168.0.2", ToElement: []Connection{{rand.Int31n(10)}, {0}, {rand.Int31n(10)}, {rand.Int31n(10)}}},
			{Name: "service 3", IP: "192.168.0.3", ToElement: []Connection{{rand.Int31n(10)}, {rand.Int31n(10)}, {0}, {rand.Int31n(10)}}},
			{Name: "service 4", IP: "192.168.0.4", ToElement: []Connection{{rand.Int31n(10)}, {0}, {rand.Int31n(10)}, {0}}},
		},
	}

	res, _ := json.Marshal(c)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "%s", res)
}
