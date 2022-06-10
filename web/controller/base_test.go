package controller

import (
	"goSearcher/searcher/db"
	"testing"
)

func TestQuery(t *testing.T) {
	db.ConnectMySql()
	relatedSearch("佟丽娅", []int{62205293, 60795789, 2222})
}
