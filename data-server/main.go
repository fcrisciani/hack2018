package main

import "github.com/fcrisciani/hack2018/data-server/methods"

func main() {
	s := methods.New()
	s.Init()
	s.Start()
	// block forever
	select {}
}
