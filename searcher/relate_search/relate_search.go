package relate_search

import (
	"fmt"
	"goSearcher/searcher/db"
	"goSearcher/searcher/model"
	"goSearcher/searcher/utils"
	"log"
	"math"
	"sort"
)

func GetRelatedSearchQueries(content string, docIds []int) (queries []string) {
	realateSearchIds := relatedSearch(content, docIds)
	queries = getQueriesByIds(realateSearchIds)
	return
}

//-------------internal functions---------------
func relatedSearch(content string, docIds []int) (queryIds []int) {
	var newqid int
	query := &model.Query{
		Query:  content,
		DocIds: fmt.Sprintf("%v", docIds),
	}
	results := db.MysqlDB.Create(query)
	if results.Error != nil {
		var queryModel = model.Query{Query: content}
		db.MysqlDB.Where("query = ?", content).First(&queryModel)
		newqid = queryModel.ID
	} else {
		newqid = query.ID
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
	querySimMap := make([]model.QuerySim, 0)
	for k, v := range similarity {
		querySimMap = append(querySimMap, model.QuerySim{QueryId: k, Similarity: v})
	}
	sort.Slice(querySimMap, func(i, j int) bool {
		return querySimMap[i].Similarity > querySimMap[j].Similarity
	})
	for i := 0; i < 10; i++ {
		if i == len(querySimMap) {
			break
		}
		queryIds = append(queryIds, querySimMap[i].QueryId)
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
