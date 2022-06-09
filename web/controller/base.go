package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goSearcher/searcher/core"
	"goSearcher/searcher/model"
	"goSearcher/searcher/utils"
	"goSearcher/searcher/words"
	"net/http"
)

func Index(c *gin.Context) {
	//userInfo := getCurrentUser(c)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"login": false})
}

func Query(c *gin.Context) {
	var documents []model.Docs
	content := c.Query("content")
	//cut content to many terms by cut model
	tokenizer := words.NewTokenizer()
	words := tokenizer.CutContent(content)
	//search in index
	for _, item := range words {
		if value, ok := core.MemoryBTree.Find(item); ok {
			docIds := utils.SplitDocIdsFromValue(fmt.Sprintf("%v", value))
			docs := core.GetDocuments(docIds)
			documents = append(documents, docs...)
		}
	}
	//continue to score
	//....
}
