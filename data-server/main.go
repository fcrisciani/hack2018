package main

import (
	"github.com/fcrisciani/hack2018/data-server/methods"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	s := methods.New()
	s.Init()
	s.Start()
	// elastic.GetPods()
	// services, _ := elastic.GetServices()
	//
	// for _, s := range services {
	// 	logrus.Infof("Pods for service:%+v", s)
	// 	pods, _ := elastic.GetPodsForService(s)
	// 	for _, p := range pods {
	// 		logrus.Infof("Connection for pod:%+v", p)
	// 		elastic.GetAllConnections(p.Status.PodIP, 0)
	// 	}
	// 	logrus.Infof("Connections for service:%+v", s)
	// 	elastic.GetAllConnections(s.Spec.ClusterIP, 0)
	// }
	// elastic.GetAllConnections("10.96.103.186", 0)
	// logrus.Infof("========")
	// elastic.GetAllConnections("192.168.195.217", 0)
	// block forever
	select {}
}
