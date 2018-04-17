package main

import (
	"github.com/fcrisciani/hack2018/data-server/elastic"
	"github.com/fcrisciani/hack2018/data-server/methods"
)

func main() {
	s := methods.New()
	s.Init()
	// s.Start()
	elastic.MatchField()
	// block forever
	// select {}
}
