package elasticsearch

import (
	"context"
	"fmt"

	"github.com/abhi/k8tcher/config"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	api_v1 "k8s.io/api/core/v1"
)

const (
	K8sResource        = "k8s.io_resource"
	K8sResourcePod     = "pod"
	K8sResourceService = "service"
)

type ElasticSearch struct {
	elasticClient *elastic.Client
}

func (e *ElasticSearch) Init(c *config.Config) error {
	logrus.Infof("Received Init :%v", c)
	client, err := elastic.NewClient(
		elastic.SetURL("http://52.42.55.249:9200"), elastic.SetSniff(false))

	if err != nil {
		logrus.Errorf("failed to initialized elastic client: %v", err)
		return fmt.Errorf("faile to initialize  elastic client: %v", err)
	}
	//esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	//if err != nil {
	// Handle error
	//	panic(err)
	//}
	//logrus.Infof("Elasticsearch version %s\n", esversion)
	e.elasticClient = client
	logrus.Infof("Done with registration!!!")
	return nil
}

func (e *ElasticSearch) ObjectCreated(obj interface{}) {
	logrus.Infof("Received object created: %v", obj)
	ctx := context.Background()
	if apiPod, ok := obj.(*api_v1.Pod); ok {
		logrus.Infof("Received Pod created: %+v", apiPod)
		_, err := e.elasticClient.Index().
			Index(K8sResource).
			Type(K8sResourcePod).
			Id(apiPod.GetName() + "-" + apiPod.GetNamespace()).
			BodyJson(apiPod).
			Do(ctx)
		if err != nil {
			// Handle error
			logrus.Errorf("failed to post pod metadata of pod %s in namespace %s to elastic search: %v",
				apiPod.GetName(), apiPod.GetNamespace(), err)
			return
		}
		return
	}
	if apiService, ok := obj.(*api_v1.Service); ok {
		logrus.Infof("Received Service created: %+v", apiService)
		//b, err := json.Marshal(apiService)
		//if err != nil {
		//logrus.Errorf("failed to parse json: %v", err)
		//}
		_, err := e.elasticClient.Index().
			Index(K8sResource).
			Type(K8sResourceService).
			Id(apiService.GetName() + "-" + apiService.GetNamespace()).
			BodyJson(apiService).
			Do(ctx)
		if err != nil {
			// Handle error
			logrus.Errorf("failed to post service metadata of service %s in namespace %s to elastic search: %v",
				apiService.GetName(), apiService.GetNamespace(), err)
			return
		}
	}
}

func (e *ElasticSearch) ObjectDeleted(obj interface{}) {
	ctx := context.Background()
	logrus.Infof("Received object deleted: %v", obj)
	if apiPod, ok := obj.(*api_v1.Pod); ok {
		logrus.Infof("Received Pod delete: %+v", apiPod)
		_, err := e.elasticClient.Delete().
			Index(K8sResource).
			Type(K8sResourcePod).
			Id(apiPod.GetName() + "-" + apiPod.GetNamespace()).
			Do(ctx)
		if err != nil {
			// Handle error
			logrus.Errorf("failed to delete pod metadata of pod %s in namespace %s to elastic search: %v",
				apiPod.GetName(), apiPod.GetNamespace(), err)
			return
		}
	}
	if apiService, ok := obj.(*api_v1.Service); ok {
		logrus.Infof("Received Service delete: %+v", apiService)
		_, err := e.elasticClient.Delete().
			Index(K8sResource).
			Type(K8sResourceService).
			Id(apiService.GetName() + "-" + apiService.GetNamespace()).
			Do(ctx)
		if err != nil {
			// Handle error
			logrus.Errorf("failed to delete service metadata of service %s in namespace %s to elastic search: %v",
				apiService.GetName(), apiService.GetNamespace(), err)
			return
		}
	}
}

func (e *ElasticSearch) ObjectUpdated(oldObj, newObj interface{}) {
	logrus.Infof("Received object updated: %v", newObj)
	if apiPod, ok := newObj.(*api_v1.Pod); ok {
		logrus.Infof("Received pod updated: %v", apiPod)
	}
}
