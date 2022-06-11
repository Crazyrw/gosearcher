package controller

import (
	"github.com/gin-gonic/gin"
	"goSearcher/searcher/core"
	"goSearcher/searcher/relate_search"
	"goSearcher/searcher/utils"
	"goSearcher/searcher/words"
	"goSearcher/web/result"
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
	//get all documents
	documents := core.GetDocuments(docIds)
	//to score: get new documents

	//relate search
	relatedSearchQueries := relate_search.GetRelatedSearchQueries(content, docIds)

	result.ResponseSuccessWithData(c, result.QueryResult{
		RelatedSearch: relatedSearchQueries,
		Documents:     documents,
	})

}
