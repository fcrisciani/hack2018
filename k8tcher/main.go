package main

import (
	"github.com/abhi/k8tcher/config"
	"github.com/abhi/k8tcher/pkg/controller"
	"github.com/abhi/k8tcher/pkg/handler"
	"github.com/abhi/k8tcher/pkg/handler/elasticsearch"
	"github.com/sirupsen/logrus"
	//"github.com/BurntSushi/toml"
)

func main() {
	//var cfg config.Config
	//if _, err := toml.DecodeFile("/k8tcher.toml", &cfg); err != nil {
	//	fmt.Printf("Failed to decode config file %v", err)
	//return
	//}
	cfg := &config.Config{
		Resource: config.Resource{
			Pod:     "yes",
			Service: "yes",
		},
		Handler: config.Handler{
			Name:   config.ElasticSearch,
			Config: config.HandlerConfig{},
		},
	}
	Run(cfg)
}

func Run(cfg *config.Config) {
	var handler handler.Handler
	switch cfg.Handler.Name {
	case config.ElasticSearch:
		handler = new(elasticsearch.ElasticSearch)
	}

	if err := handler.Init(cfg); err != nil {
		logrus.Errorf("Failed to initialize:%v", err)
		return
	}

	controller.Start(cfg, handler)
}
