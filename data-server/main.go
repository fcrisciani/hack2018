package main

import (
	"github.com/fcrisciani/hack2018/data-server/elastic"
	"github.com/fcrisciani/hack2018/data-server/methods"
)

func main() {
	s := methods.New()
	s.Init()
	// s.Start()
	elastic.GetServices()
	elastic.GetConnections("10.96.0.1", 0)
	// block forever
	// select {}
}
