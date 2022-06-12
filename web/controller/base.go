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

var page *utils.Paging
var data result.QueryResult
var relatedSearchQueries []string

func Index(c *gin.Context) {
	//userInfo := getCurrentUser(c)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func Paging(total int) *utils.Paging {

	paging := utils.CreatePaging(1, 10, total)

	return paging
}

func Query(c *gin.Context) {

	docIdsMap := make(map[int]int, 0)
	var docIds []int //union set, content's results
	var lens []int   //terms-len(docIds)

	content := c.Query("content")
	exclude := c.Query("exclued")

	fmt.Println(exclude)

	pageNum := 1

	//cut content to many terms by cut model
	tokenizer := words.NewTokenizer()
	words := tokenizer.CutContent(content)
	fmt.Println(words)
	//search in index
	for _, item := range words {
		//之后如果进行优化的话 可以并发的读
		value := core.SkipList.Search([]byte(item)).Value
		ids := utils.SplitDocIdsFromValue(string(value))
		// fmt.Println(ids)
		lens = append(lens, len(ids))
		//union docIds
		for _, id := range ids {
			_, ok := docIdsMap[id]
			if !ok {
				docIdsMap[id] = 1
			}
		}
	}
	for key := range docIdsMap {
		docIds = append(docIds, key)
	}
	// fmt.Println(docIds)
	if len(docIds) == 0 {
		result.Error("no results")
	}
	//to score: get new docIds
	rankDocuments := rank.Rank(docIds, words, lens)

	// fmt.Println("---------------------------------------")
	// fmt.Println(rankDocuments)
	//relate search

	relatedSearchQueries = relate_search.GetRelatedSearchQueries(content, docIds)
	data = result.QueryResult{
		RelatedSearch: relatedSearchQueries,
		Documents:     rankDocuments,
	}

	resultsNum := len(data.Documents)

	page = Paging(resultsNum)
	page.Page = pageNum

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageNum >= int(page.PageCount) {
		pageNum = int(page.PageCount)
	}

	// 分页返回数据
	var pageData = result.QueryResult{
		RelatedSearch: relatedSearchQueries,
		Documents:     data.Documents[10*(pageNum-1) : 10*pageNum],
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"content": content, "page": page, "Data": pageData})

}

func GetLastPage(c *gin.Context) {
	pageNum := page.Page - 1

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageNum >= page.PageCount {
		pageNum = page.PageCount
	}

	// 分页返回数据
	var pageData = result.QueryResult{
		RelatedSearch: relatedSearchQueries,
		Documents:     data.Documents[10*(pageNum-1) : 10*pageNum],
	}
	page.Page = pageNum
	content := c.Query("content")

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"content": content, "page": page, "Data": pageData})
}

func GetNextPage(c *gin.Context) {
	pageNum := page.Page + 1

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageNum >= int(page.PageCount) {
		pageNum = int(page.PageCount)
	}

	// 分页返回数据
	var pageData = result.QueryResult{
		RelatedSearch: relatedSearchQueries,
		Documents:     data.Documents[10*(pageNum-1) : 10*pageNum],
	}
	page.Page = pageNum

	content := c.Query("content")

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"content": content, "page": page, "Data": pageData})
}
