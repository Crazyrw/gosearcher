package core

import (
	"goSearcher/searcher/btree"
)

var MemoryBTree *btree.BPlusTree

func Initialize() {
	//create memory btree
	MemoryBTree = CreateMemoryBtree("searcher/data/terms/dictionary1.txt")
}
