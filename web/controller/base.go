package controller

import (
	"fmt"
	"goSearcher/searcher/core"
	"goSearcher/searcher/model"
	"goSearcher/searcher/rank"
	"goSearcher/searcher/relate_search"
	"goSearcher/searcher/utils"
	"goSearcher/searcher/words"
	"goSearcher/web/result"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var page *utils.Paging
var data result.QueryResult
var relatedSearchQueries []string
var rankDocIds []int //union set, content's results
var terms []string

type Position struct {
	Start int
	End   int
}
type DocumentPos struct {
	Document model.Docs
	Pos      []Position
}

func Index(c *gin.Context) {
	//userInfo := getCurrentUser(c)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func Paging(total int) *utils.Paging {

	paging := utils.CreatePaging(1, 10, total)

	return paging
}

func Query(c *gin.Context) {
	rankDocIds = nil
	terms = nil
	relatedSearchQueries = nil

	docIdsMap := make(map[int]int, 0)
	var lens []int //terms-len(docIds)

	content := c.Query("content")
	exclude := c.Query("exclued")

	excludeDocIds := queryExclude(exclude)

	pageNum := 1

	//cut content to many terms by cut model
	tokenizer := words.NewTokenizer()
	terms = tokenizer.CutContent(content)
	if len(terms) == 0 {
		c.HTML(http.StatusBadRequest, "index.tmpl", gin.H{"State": false})
		return
	}
	//search in index
	for _, item := range terms {
		//之后如果进行优化的话 可以并发的读
		value := core.SkipList.Search([]byte(item)).Value
		ids := utils.SplitDocIdsFromValue(string(value))
		// fmt.Println(ids)
		lens = append(lens, len(ids))
		//union docIds
		for _, id := range ids {
			_, ok := docIdsMap[id]
			_, ok1 := excludeDocIds[id]
			if !ok && !ok1 {
				docIdsMap[id] = 1
			}
		}
	}
	var docIds []int
	for key := range docIdsMap {
		docIds = append(docIds, key)
	}
	// fmt.Println(docIds)
	if len(docIds) == 0 {
		c.HTML(http.StatusBadRequest, "index.tmpl", gin.H{"State": false})
		return
	}
	//to score: get new docIds
	rankDocIds = rank.Rank(docIds, terms, lens)

	// fmt.Println("---------------------------------------")
	// fmt.Println(rankDocuments)
	//relate search

	relatedSearchQueries = relate_search.GetRelatedSearchQueries(content, docIds)

	resultsNum := len(rankDocIds)

	page = Paging(resultsNum)
	page.Page = pageNum

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageNum >= int(page.PageCount) {
		pageNum = int(page.PageCount)
	}

	//第一页的数据 返回
	firstDocIds := rankDocIds[10*(pageNum-1) : 10*pageNum]

	docs := utils.GetDocumentsFor(firstDocIds)

	finalDocs := hightLight(terms, docs)
	// 分页返回数据
	var pageData = result.QueryResult{
		RelatedSearch: relatedSearchQueries,
		Documents:     finalDocs,
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

	page.Page = pageNum
	content := c.Query("content")

	docIds := rankDocIds[10*(pageNum-1) : 10*pageNum]
	docs := utils.GetDocumentsFor(docIds)
	finalDocs := hightLight(terms, docs)
	// 分页返回数据
	var pageData = result.QueryResult{
		RelatedSearch: relatedSearchQueries,
		Documents:     finalDocs,
	}

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

	page.Page = pageNum
	content := c.Query("content")

	docIds := rankDocIds[10*(pageNum-1) : 10*pageNum]
	docs := utils.GetDocumentsFor(docIds)
	finalDocs := hightLight(terms, docs)
	// 分页返回数据
	var pageData = result.QueryResult{
		RelatedSearch: relatedSearchQueries,
		Documents:     finalDocs,
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"content": content, "page": page, "Data": pageData})
}

func queryExclude(exclude string) map[int]int {
	tokenizer := words.NewTokenizer()
	excludeTerms := tokenizer.CutContent(exclude)
	docIdsMap := make(map[int]int, 0)
	if len(excludeTerms) == 0 {
		fmt.Println("excludeTerms == 0")
		result.Error("no results")
		return nil
	}
	//search in index
	for _, item := range excludeTerms {
		//之后如果进行优化的话 可以并发的读
		value := core.SkipList.Search([]byte(item)).Value
		ids := utils.SplitDocIdsFromValue(string(value))
		//union docIds
		for _, id := range ids {
			_, ok := docIdsMap[id]
			if !ok {
				docIdsMap[id] = 1
			}
		}
	}
	return docIdsMap
}

func hightLight(terms []string, documents []model.Docs) []model.Docs {
	var documentPos []DocumentPos
	for _, item := range documents {
		content := item.Caption
		var pos []Position
		for _, term := range terms {
			start := strings.Index(content, term)
			if start == -1 {
				continue
			}
			end := start + len(term) - 1
			p := Position{
				Start: start,
				End:   end,
			}
			pos = append(pos, p)
		}
		if pos == nil {
			continue
		}
		documentPosItem := DocumentPos{
			Document: item,
			Pos:      pos,
		}
		// fmt.Println(documentPosItem)
		documentPos = append(documentPos, documentPosItem)
	}

	// fmt.Println(documentPos)
	spanStart := "<span style=\"color:red\">"
	spanEnd := "</span>"
	finalDocs := make([]model.Docs, 0)
	for _, docPos := range documentPos {
		caption := docPos.Document.Caption
		for _, pos := range docPos.Pos {
			start := pos.Start
			end := pos.End
			word := caption[start : end+1]
			docPos.Document.Caption = strings.ReplaceAll(docPos.Document.Caption, word, spanStart+word+spanEnd)
		}
		finalDocs = append(finalDocs, docPos.Document)
	}
	// return documentPos
	return finalDocs
}
