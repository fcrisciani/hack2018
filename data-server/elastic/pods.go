package elastic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

type podMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
}

type podStatus struct {
	HostIP string `json:"hostIP"`
	PodIP  string `json:"podIP"`
	Phase  string `json:"phase"`
}

type Pod struct {
	Meta   podMetadata `json:"metadata"`
	Status podStatus   `json:"status"`
}

func GetPods(c *elastic.Client) ([]*Pod, error) {

	// termQuery := elastic.NewTermQuery("dest_ip", "172.31.39.84")
	searchResult, err := c.Search().
		Index("k8s.io_resource"). // search in index "twitter"
		Type("pod").
		Size(10000).
		// Query(termQuery).             // specify the query
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
	logrus.Debugf("Query took %d milliseconds", searchResult.TookInMillis)
	return parsePods(searchResult)
}

func GetPodsForService(srv *Service) ([]*Pod, error) {
	c, err := NewClient("52.42.55.249", "9200")
	if err != nil {
		panic(err)
	}

	if len(srv.Meta.Labels) == 0 {
		return nil, nil
	}

	labelToMatch := make([]elastic.Query, 0, len(srv.Meta.Labels))
	for k, v := range srv.Meta.Labels {
		termQuery := elastic.NewMatchPhraseQuery(fmt.Sprintf("metadata.labels.%s", k), v)
		labelToMatch = append(labelToMatch, termQuery)
	}
	shouldClause := elastic.NewBoolQuery().Filter(elastic.NewBoolQuery().Should(labelToMatch...))

	// logrus.Infof("Query is:%+v", shouldClause)

	searchResult, err := c.Search().
		Index("k8s.io_resource"). // search in index "twitter"
		Type("pod").
		Query(shouldClause). // specify the query
		Size(10000).
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
	logrus.Debugf("Query took %d milliseconds", searchResult.TookInMillis)
	return parsePods(searchResult)
}

func parsePods(searchResult *elastic.SearchResult) ([]*Pod, error) {
	podList := make([]*Pod, 0, searchResult.Hits.TotalHits)
	if searchResult.Hits.TotalHits > 0 {
		logrus.Debugf("Found a total of %d pod, len hits:%d", searchResult.Hits.TotalHits, len(searchResult.Hits.Hits))

		// Iterate through results
		for i, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Pod
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			logrus.Debugf("%d) Pod %+v", i, t)
			podList = append(podList, &t)
		}
	} else {
		// No hits
		logrus.Debugf("Found no pod")
		return nil, nil
	}

	return podList, nil
}
