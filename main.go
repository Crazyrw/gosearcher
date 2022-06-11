package main

import (
	"goSearcher/searcher/db"
	"goSearcher/web/router"
)

//var stat runtime.MemStats
func init() {
	db.ConnectMySql()
}
func main() {
	//searcher service
	//core.Initialize()
	//web service
	router := router.SetupRouter()
	router.Run(":9090")
}
