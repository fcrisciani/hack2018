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
	Protocol int    `json:"orig.ip.protocol"`
	SrcPort  int    `json:"orig.l4.sport"`
	DstPort  int    `json:"orig.l4.dport"`
}

func GetAllConnections(IP string, protocol int) ([]*Connection, error) {
	c, err := NewClient("52.42.55.249", "9200")
	if err != nil {
		panic(err)
	}

	shouldClause := elastic.NewBoolQuery().Should(elastic.NewMatchPhraseQuery("dest_ip", IP), elastic.NewMatchPhraseQuery("src_ip", IP))
	if protocol != 0 {
		shouldClause = shouldClause.Must(elastic.NewMatchPhraseQuery("ip.protocol", protocol))
	}
	searchResult, err := c.Search().
		Index("logstash-*"). // search in index "twitter"
		Query(shouldClause). // specify the query
		// Sort("Time", true).           // sort by "user" field, ascending
		// From(0).Size(10).        // take documents 0-9
		// Pretty(true).            // pretty print request and response JSON
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
			fmt.Printf("Connection %+v\n", t)

			connections = append(connections, &t)
		}
	} else {
		// No hits
		fmt.Print("Found no connections\n")
		return nil, nil
	}

	return connections, nil
}
