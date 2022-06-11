package core

import (
	"goSearcher/searcher/skip_list"
)

//var MemoryBTree *btree.BPlusTree
var SkipList *skip_list.SkipList

func Initialize() {
	//create memory btree
	//MemoryBTree = CreateMemoryBtree("searcher/data/terms/dictionary1.txt")

	//create memory skipList
	SkipList = CreateSkipList("searcher/data/terms/dictionary1.txt")
}
