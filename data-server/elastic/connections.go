package elastic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"
)

func NewClient(ip, port string) (*elastic.Client, error) {
	// elastic.SetSniff(false)
	return elastic.NewClient(elastic.SetURL(fmt.Sprintf("http://%s:%s/", ip, port)), elastic.SetSniff(false))
}

type Connection struct {
	SrcIP    string `json:"src_ip"`
	DstIP    string `json:"dest_ip"`
	Protocol string `json:"oob.prefix"`
	SrcPort  int    `json:"src_port"`
	DstPort  int    `json:"dest_port"`
}

func GetConnections(IP string, protocol int) ([]*Connection, error) {
	c, err := NewClient("52.42.55.249", "9200")
	if err != nil {
		panic(err)
	}

	// termQuery := elastic.NewTermQuery("dest_ip", "172.31.39.84")
	searchResult, err := c.Search().
		Index("k8s.io_resource"). // search in index "twitter"
		Type("service").
		// Query(termQuery).             // specify the query
		// Sort("Time", true).           // sort by "user" field, ascending
		From(0).Size(10).        // take documents 0-9
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	shouldClause := elastic.NewBoolQuery().Should(elastic.NewTermQuery("dest_ip", IP), elastic.NewTermQuery("src_ip", IP))
	if protocol != 0 {
		shouldClause = shouldClause.Must(elastic.NewTermQuery("ip.protocol", protocol))
	}
	searchResult, err = c.Search().
		Index("logstash-2018.04.17"). // search in index "twitter"
		Query(shouldClause).          // specify the query
		// Sort("Time", true).           // sort by "user" field, ascending
		From(0).Size(10).        // take documents 0-9
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		// Handle error
		panic(err)
	}
	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	connections := make([]*Connection, 0, searchResult.Hits.TotalHits)
	// Number of hits
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d packets\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Connection
			// fmt.Printf("%s", *hit.Source)
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Packet %+v\n", t)

			connections = append(connections, &t)
		}
	} else {
		// No hits
		fmt.Print("Found no connections\n")
		return nil, nil
	}

	return connections, nil
}
