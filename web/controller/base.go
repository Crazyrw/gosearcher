package controller

import (
	"fmt"
	"goSearcher/searcher/core"
	"goSearcher/searcher/rank"
	"goSearcher/searcher/relate_search"
	"goSearcher/searcher/utils"
	"goSearcher/searcher/words"
	"goSearcher/web/result"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	//userInfo := getCurrentUser(c)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}
func Query(c *gin.Context) {

	docIdsMap := make(map[int]int)
	var docIds []int //union set, content's results
	var lens []int   //term-len(docIds)
	content := c.Query("content")
	//cut content to many terms by cut model
	tokenizer := words.NewTokenizer()
	words := tokenizer.CutContent(content)
	fmt.Println(words)
	//search in index
	for _, item := range words {
		//之后如果进行优化的话 可以并发的读
		value := core.SkipList.Search([]byte(item)).Value

		ids := utils.SplitDocIdsFromValue(string(value))
		lens = append(lens, len(ids))
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
	rankDocuments := rank.Rank(docIds, words, lens)
	// fmt.Println("---------------------------------------")
	// fmt.Println(rankDocuments)
	//relate search
	relatedSearchQueries := relate_search.GetRelatedSearchQueries(content, docIds)
	var d = result.QueryResult{
		RelatedSearch: relatedSearchQueries,
		Documents:     rankDocuments,
	}
	//fmt.Println(d.Documents[1], d.Documents[2])
	result.ResponseSuccessWithData(c, d)

}
