package elastic

import (
	"context"
	"encoding/json"
	"fmt"
)

type metadata struct {
	ServiceName string `json:"name"`
}

type serviceSpec struct {
	ClusterIP string `json:"clusterIP"`
}

type Service struct {
	Meta metadata `json:"metadata"`
	Spec serviceSpec
}

func GetServices() ([]*Service, error) {
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

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	serviceList := make([]*Service, 0, searchResult.Hits.TotalHits)
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d packets\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Service
			// fmt.Printf("%s", *hit.Source)
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			// fmt.Printf("Service %s\n", *hit.Source)

			fmt.Printf("Service %+v\n", t)
			serviceList = append(serviceList, &t)
		}
	} else {
		// No hits
		fmt.Print("Found no packets\n")
		return nil, nil
	}

	return serviceList, nil
}
