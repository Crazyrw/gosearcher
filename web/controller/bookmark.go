package controller

import (
	"goSearcher/searcher/db"
	"goSearcher/web/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

//创建用户书签
func Create_bookmark(c *gin.Context) {

	var user model.User

	phone := c.Query("phone")
	bookmark_name := c.Query("bookmark_name")
	db.MysqlDB.First(&user, "phone = ?", phone)

	if err := db.MysqlDB.Model(&user).Update("bookmark_name", bookmark_name).Error; err != nil {
		//返回结果
		c.JSON(http.StatusBadRequest, gin.H{"message": "创建失败"})
		return
	}

	//返回结果
	c.JSON(http.StatusOK, gin.H{"message": "成功"})

}

//添加书签
func Add_bookmark(c *gin.Context) {

	phone := c.Query("phone")
	docid := c.Query("docid")
	var docs model.Docs
	db.MysqlDB.Where("id = ?", docid).First(&docs)
	caption := docs.Caption
	newBookmark := model.Bookmark{
		Phone:   phone,
		DocId:   docid,
		Caption: caption,
	}

	db.MysqlDB.Create(&newBookmark)

	// 返回结果
	c.JSON(http.StatusOK, gin.H{"message": "OK"})

}

//删除单个书签
func Delete_bookmark(c *gin.Context) {

	var bookmark model.Bookmark

	phone := c.Query("phone")
	docid := c.Query("docid")

	db.MysqlDB.Unscoped().Where("phone = ? and doc_id = ? ", phone, docid).Delete(&bookmark)

	//返回结果
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})

}

//删除所有书签
func DeleteAll_bookmark(c *gin.Context) {

	var bookmark model.Bookmark
	var user model.User

	phone := c.Query("phone")

	db.MysqlDB.First(&user, "phone = ?", phone)
	db.MysqlDB.Model(&user).Update("bookmark_name", "")

	db.MysqlDB.Unscoped().Where("phone = ? ", phone).Delete(&bookmark)

	//返回结果
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})

}

//获得用户收藏的书签内容
func Get_bookmark(c *gin.Context) {

	phone := c.Query("phone")

	var bookmarks = []model.Bookmark{}

	db.MysqlDB.Where("phone = ?", phone).Find(&bookmarks)

	c.JSON(http.StatusOK, gin.H{"bookmarks": bookmarks})

}
