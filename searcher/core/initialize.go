package core

import (
	"goSearcher/searcher/btree"
	"goSearcher/searcher/db"
)

var MemoryBTree *btree.BPlusTree

func init() {
	db.ConnectMySql()
}
func Initialize() {
	MemoryBTree = CreateMemoryBtree("searcher/data/terms/dictionary.txt")
}
