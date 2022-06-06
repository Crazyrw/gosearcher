package btree

import (
	"fmt"
	"testing"
)

func TestNewTree(t *testing.T) {
	tree, errsr := NewTree("../../searcher/data/index/data.db")
	if errsr != nil {
		panic(errsr)
	}
	defer tree.Close()
}
func TestInsert(t *testing.T) {
	tree, errsr := NewTree("../../searcher/data/index/data.db")
	if errsr != nil {
		panic(errsr)
	}
	defer tree.Close()
	errsr = tree.Insert("佟丽娅", "[251917]")
	if errsr != nil {
		panic(errsr)
	}
}
func TestFind(t *testing.T) {
	tree, errsr := NewTree("../../searcher/data/index/data.db")
	if errsr != nil {
		panic(errsr)
	}
	defer tree.Close()
	values, errsr := tree.Find("佟丽娅")
	if errsr != nil {
		panic(errsr)
	}
	fmt.Println(values)
}
