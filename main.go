package main

import (
	"goSearcher/searcher/core"
	"goSearcher/web/router"
)

//var stat runtime.MemStats

func main() {
	//searcher service
	core.Initialize()
	//web service
	router := router.SetupRouter()
	router.Run(":9090")
}
