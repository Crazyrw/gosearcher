package rank

import (
	"fmt"
	"goSearcher/searcher/db"
	"goSearcher/searcher/model"
	"testing"
)

func TestRank(t *testing.T) {
	db.ConnectMySql()
	docRank := Rank([]int{89149531, 81094679, 83493790, 92311417, 88362644, 85460045, 813089149531, 7803, 80783350, 89806811, 92024184, 87369755, 82968689, 83644567, 83922887, 87451911, 80831224, 91807837}, []string{"golang"}, []int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5})
	fmt.Println(docRank)
}
func TestCount(t *testing.T) {
	db.ConnectMySql()
	var count int64
	db.MysqlDB.Model(&model.Docs{}).Count(&count)
	fmt.Println(count)
}
