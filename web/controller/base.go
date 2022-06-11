package controller

import (
	"github.com/gin-gonic/gin"
	"goSearcher/searcher/core"
	"goSearcher/searcher/relate_search"
	"goSearcher/searcher/utils"
	"goSearcher/searcher/words"
	"goSearcher/web/result"
	"math"
	"net/http"
)

func Index(c *gin.Context) {
	//userInfo := getCurrentUser(c)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}
func Query(c *gin.Context) {

	docIdsMap := make(map[int]int)
	var docIds []int //union set, content's results
	content := c.Query("content")
	//cut content to many terms by cut model
	tokenizer := words.NewTokenizer()
	words := tokenizer.CutContent(content)
	//search in index
	for _, item := range words {
		//之后如果进行优化的话 可以并发的读
		value := core.SkipList.Search([]byte(item)).Value

		ids := utils.SplitDocIdsFromValue(string(value))
		//union docIds
		for _, id := range ids {
			_, ok := docIdsMap[id]
			if !ok {
				docIdsMap[id] = 1
				docIds = append(docIds, id)
			}
		}
	}
	if len(docIds) == 0 {
		result.Error("no results")
	}
	//to score: get new docIds

	//relate search
	relatedSearchQueries := relate_search.GetRelatedSearchQueries(content, docIds)

	result.ResponseSuccessWithData(c, result.QueryResult{
		RelatedSearch: relatedSearchQueries,
		Documents:     docIds,
	})

}

func rank(query string, docs []string) (docsScore []float64) {
	tokenizer := words.NewTokenizer()
	queryWords := tokenizer.CutContent(query)         //query的分词结果
	IDFs := make(map[string]float64, len(queryWords)) //保存query中每个word的idf值
	for _, word := range queryWords {
		var num float64 = 0
		for _, doc := range docs {
			docWords := tokenizer.CutContent(doc) //doc的分词结果
			flag := Find(docWords, word)
			if flag == true {
				num++
			}
		}
		//df := num / float64(len(docs))
		idf := math.Log2(float64(len(docs)) / num)
		IDFs[word] = idf
	}

	docsScore = make([]float64, len(docs), len(docs))
	for k, doc := range docs {
		score := getScore(query, doc, IDFs)
		docsScore[k] = score
	}
	return docsScore
}

//获取query对于某个doc的tf-idf值
func getScore(query string, doc string, IDFs map[string]float64) (score float64) {
	score = 0
	//var DF float64 = 0
	tokenizer := words.NewTokenizer()
	queryWords := tokenizer.CutContent(query) //query的分词结果
	docWords := tokenizer.CutDoc(doc)         //doc的分词结果
	for _, word := range queryWords {
		count := Count(docWords, word)
		var tf float64 = count / float64(len(docWords))
		idf := IDFs[word]
		score += tf * idf
	}

	return score
}

//判断string切片中含某个string的个数
func Count(slice []string, val string) (count float64) {
	count = 0
	for _, item := range slice {
		if item == val {
			count++
		}
	}
	return count
}

//判断string切片中是否含某个string
func Find(slice []string, val string) (flag bool) {
	flag = false
	for _, item := range slice {
		if item == val {
			flag = true
			break
		}
	}
	return flag
}
