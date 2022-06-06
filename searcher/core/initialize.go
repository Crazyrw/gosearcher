package core

import (
	"goSearcher/searcher/db"
	"sync"
)

var kv sync.Map

func init() {
	db.ConnectMySql()
}
func Initialize() {
	//createInvertIndex() //落盘
}
