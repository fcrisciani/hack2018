package methods

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/fcrisciani/hack2018/data-server/elastic"
	"github.com/sirupsen/logrus"
)

type Connection struct {
	Total int32                `json:"total"`
	Flows []elastic.Connection `json:"flows"`
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
			{Name: "service 1", IP: "192.168.0.1", ToElement: []Connection{{0, nil}, {rand.Int31n(10), nil}, {rand.Int31n(10), nil}, {rand.Int31n(10), nil}}},
			{Name: "service 2", IP: "192.168.0.2", ToElement: []Connection{{rand.Int31n(10), nil}, {0, nil}, {rand.Int31n(10), nil}, {rand.Int31n(10), nil}}},
			{Name: "service 3", IP: "192.168.0.3", ToElement: []Connection{{rand.Int31n(10), nil}, {rand.Int31n(10), nil}, {0, nil}, {rand.Int31n(10), nil}}},
			{Name: "service 4", IP: "192.168.0.4", ToElement: []Connection{{rand.Int31n(10), nil}, {0, nil}, {rand.Int31n(10), nil}, {0, nil}}},
		},
	}

	res, _ := json.Marshal(c)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "%s", res)
}
