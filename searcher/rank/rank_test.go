package rank

import (
	"fmt"
	"goSearcher/searcher/db"
	"goSearcher/searcher/model"
	"testing"
)

func TestRank(t *testing.T) {
	db.ConnectMySql()
	docRank := Rank([]int{80553124, 81607254, 80564911, 81611175}, []string{"字节", "跳动"}, []int{5, 5, 5, 5})
	fmt.Println(docRank)
}
func TestCount(t *testing.T) {
	db.ConnectMySql()
	var count int64
	db.MysqlDB.Model(&model.Docs{}).Count(&count)
	fmt.Println(count)
}
