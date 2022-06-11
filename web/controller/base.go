package controller

import (
	"fmt"
	"goSearcher/searcher/core"
	"goSearcher/searcher/db"
	"goSearcher/searcher/model"
	"goSearcher/searcher/utils"
	"goSearcher/searcher/words"
	"goSearcher/web/result"
	"log"
	"math"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	//userInfo := getCurrentUser(c)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}
func Query(c *gin.Context) {

	docIdsMap := make(map[int]int)
	var docIds []int //union set
	content := c.Query("content")
	//cut content to many terms by cut model
	tokenizer := words.NewTokenizer()
	words := tokenizer.CutContent(content)
	fmt.Println(words)
	//search in index
	for _, item := range words {
		//之后如果进行优化的话 可以并发的读
		if value, ok := core.MemoryBTree.Find(item); ok {
			ids := utils.SplitDocIdsFromValue(fmt.Sprintf("%v", value))
			//union docIds
			for _, id := range ids {
				_, ok := docIdsMap[id]
				if !ok {
					docIdsMap[id] = 1
					docIds = append(docIds, id)
				}
			}
		}
	}
	//relate search
	realateSearchIds := relatedSearch(content, docIds)
	queries := getQueriesByIds(realateSearchIds)
	result.ResponseSuccessWithData(c, queries)

	//continue to score

}

//-------------internal functions---------------
// relate search
type querySim struct {
	queryId    int
	similarity float64
}

func relatedSearch(content string, docIds []int) (queryIds []int) {

	var queryModel = model.Query{Query: content}
	results := db.MysqlDB.First(&queryModel)
	var newqid int
	if results.Error == nil {
		// save to mysql-query
		query := &model.Query{
			Query:  content,
			DocIds: fmt.Sprintf("%v", docIds),
		}
		db.MysqlDB.Create(query)
		newqid = query.ID
	} else {
		newqid = queryModel.ID
	}

	//get add queries
	var allQueries []model.Query
	result := db.MysqlDB.Find(&allQueries)
	if result.Error != nil {
		log.Println("get all queries error ", result.Error)
	}
	queryDocIdMap := make(map[int][]int)
	for _, item := range allQueries {
		queryDocIdMap[item.ID] = utils.SplitDocIdsFromValue(item.DocIds)
	}
	//invert queryid-docids -> docid-queryids
	docIdQueryMap := preProcess(queryDocIdMap)
	//intersection newquery & otherqueries
	interSection := queryInterSection(docIdQueryMap, newqid)
	//similarity newquery & otherqueries
	similarity := querySimMatrix(queryDocIdMap, newqid, interSection)
	//map转成结构切片按照Sim进行降序排序
	querySimMap := make([]querySim, 0)
	for k, v := range similarity {
		querySimMap = append(querySimMap, querySim{queryId: k, similarity: v})
	}
	sort.Slice(querySimMap, func(i, j int) bool {
		return querySimMap[i].similarity > querySimMap[j].similarity
	})
	for i := 0; i < 10; i++ {
		if i == len(querySimMap) {
			break
		}
		queryIds = append(queryIds, querySimMap[i].queryId)
	}
	return
}
func preProcess(qmap map[int][]int) map[int][]int {
	// dict结构如下：
	//     {"Query1": {DocID1, DocID2, DocID3,...}
	//      "Query2": {DocID12, DocID5, DocID8,...}
	//      ...
	//     }
	// 建立倒排表
	doc_q := make(map[int][]int)
	for qid, docids := range qmap {
		for _, docid := range docids {
			doc_q[docid] = append([]int{qid}, doc_q[docid]...)
		}
	}
	return doc_q
}
func queryInterSection(doc_q map[int][]int, newqid int) map[int]int {
	//建立newquery与其他query的doc交集dict, 其中C[v]代表的含义是查询newqid和查询v之间共同doc数
	queryInterSection := make(map[int]int)
	for _, qids := range doc_q {
		if in(newqid, qids) {
			for _, qid := range qids {
				if qid == newqid {
					continue
				}
				queryInterSection[qid] += 1
			}
		}
	}
	return queryInterSection
}
func in(id int, ids []int) bool {
	for _, v := range ids {
		if id == v {
			return true
		}
	}
	return false
}
func querySimMatrix(qmap map[int][]int, newqid int, queryInterSection map[int]int) map[int]float64 {
	//建立query相似度矩阵
	querySimMatrix := make(map[int]float64)
	nnew := len(qmap[newqid])
	for qid, nums := range queryInterSection {
		nelse := len(qmap[qid])
		querySimMatrix[qid] = float64(nums) / math.Sqrt(float64(nnew*nelse))
	}
	return querySimMatrix
}
func getQueriesByIds(queryIds []int) (queries []string) {
	var queriesQ []model.Query
	result := db.MysqlDB.Find(&queriesQ, queryIds)
	if result.Error != nil {
		log.Println("getQueriesByIds err ", result.Error)
	}
	for _, item := range queriesQ {
		queries = append(queries, item.Query)
	}
	return
}
