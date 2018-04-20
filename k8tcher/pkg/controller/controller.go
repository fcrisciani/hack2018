package controller

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhi/k8tcher/config"
	"github.com/abhi/k8tcher/pkg/handler"

	"github.com/sirupsen/logrus"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	core_v1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	POD     = "pod"
	SERVICE = "service"
)

const (
	maxRetries = 5
)

var serverStartTime time.Time

// Controller holds the k8s controller information
type Controller struct {
	kubeclient kubernetes.Interface
	queue      workqueue.RateLimitingInterface
	informer   cache.SharedIndexInformer
	handler    handler.Handler
}

func Start(conf *config.Config, handler handler.Handler) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Errorf("failed to initialize in-cluster config: %v", err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Errorf("failed to create api server client: %v", err)
	}

	if conf.Resource.Pod != "" {
		informer := core_v1.NewPodInformer(clientset, meta_v1.NamespaceAll, 0, cache.Indexers{})
		c := newController(clientset, handler, informer, POD)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Service != "" {
		informer := core_v1.NewServiceInformer(clientset, meta_v1.NamespaceAll, 0, cache.Indexers{})

		c := newController(clientset, handler, informer, SERVICE)
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func newController(client kubernetes.Interface, handler handler.Handler, informer cache.SharedIndexInformer, resource string) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			logrus.WithField("event", resource).Infof("Received Add event for %v: %s", resource, key)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			logrus.WithField("event", resource).Infof("Received Update event for %v: %s", resource, key)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			logrus.WithField("event", resource).Infof("Received Delete event for %v: %s", resource, key)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	return &Controller{
		kubeclient: client,
		informer:   informer,
		queue:      queue,
		handler:    handler,
	}
}

// Run starts the kubewatch controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	logrus.Info("Starting k8tcher controller")
	serverStartTime = time.Now().Local()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	logrus.Info("k8tcher controller synced and ready")

	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced is required for the cache.Controller interface.
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// LastSyncResourceVersion is required for the cache.Controller interface.
func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	logrus.Infof("processNextItem: %+v", key)
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.processItem(key.(string))
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(key)
	} else if c.queue.NumRequeues(key) < maxRetries {
		logrus.Errorf("Error processing %s (will retry): %v", key, err)
		c.queue.AddRateLimited(key)
	} else {
		// err != nil and too many retries
		logrus.Errorf("Error processing %s (giving up): %v", key, err)
		c.queue.Forget(key)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processItem(key string) error {
	logrus.Infof("processItem: %v", key)
	obj, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", key, err)
	}
	// get object's metedata
	//objectMeta := getObjectMetaData(obj)

	if !exists {
		c.handler.ObjectDeleted(obj)
		return nil
	}
	logrus.Infof("")
	// compare CreationTimestamp and serverStartTime and alert only on latest events
	// Could be Replaced by using Delta or DeltaFIFO
	//if objectMeta.CreationTimestamp.Sub(serverStartTime).Seconds() > 0 {
	c.handler.ObjectCreated(obj)
	//}
	return nil
}

func getObjectMetaData(obj interface{}) meta_v1.ObjectMeta {

	var objectMeta meta_v1.ObjectMeta

	switch object := obj.(type) {
	case *api_v1.Service:
		objectMeta = object.ObjectMeta
	case *api_v1.Pod:
		objectMeta = object.ObjectMeta
	}
	return objectMeta
}
