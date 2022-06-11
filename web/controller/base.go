package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	//userInfo := getCurrentUser(c)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func Query(c *gin.Context) {
	//var documents []model.Docs
	//content := c.Query("content")
	////cut content to many terms by cut model
	//tokenizer := words.NewTokenizer()
	//words := tokenizer.CutContent(content)
	//search in index
	//for _, item := range words {
	//	if value, ok := core.MemoryBTree.Find(item); ok {
	//		docIds := utils.SplitDocIdsFromValue(fmt.Sprintf("%v", value))
	//		docs := core.GetDocuments(docIds)
	//		documents = append(documents, docs...)
	//	}
	//}
	//continue to score
	//....

	// 模拟数据
	resp := gin.H{
		"url":   "https://go.dev/",
		"title": "The Go Programming Language",
		"img":   "https://www.google.com/imgres?imgurl=https%3A%2F%2Fpbs.twimg.com%2Fprofile_images%2F1142154201444823041%2FO6AczwfV_400x400.png&imgrefurl=https%3A%2F%2Ftwitter.com%2Fgolang&tbnid=AfPvlfaD_83mKM&vet=12ahUKEwjhu8OG-qH4AhWJS_UHHf6pAG8QMygAegUIARCGAQ..i&docid=KnPnz5GJJTt4DM&w=400&h=400&q=go&ved=2ahUKEwjhu8OG-qH4AhWJS_UHHf6pAG8QMygAegUIARCGAQ",
	}

	c.HTML(http.StatusOK, "index.tmpl", resp)
}
