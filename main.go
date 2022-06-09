package main

import (
	"goSearcher/web/router"
)

//var stat runtime.MemStats

func main() {
	//searcher service
	//core.Initialize()

	// DB Init
	//model.UserDBInit()

	//web service
	router := router.SetupRouter()
	router.Run(":9090")
}
