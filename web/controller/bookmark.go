package controller

import (
	"goSearcher/searcher/db"
	"goSearcher/web/model"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func setCurrentBookmark(c *gin.Context, bookmark []model.Bookmark) {
	session := sessions.Default(c)
	session.Set("currentBookmark", bookmark)
	// 一定要Save否则不生效，若未使用gob注册User结构体，调用Save时会返回一个Error
	session.Save()
}

func delCurrentBookmark(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("currentBookmark")
	session.Save()
}

//创建用户书签
func Create_bookmark(c *gin.Context) {

	var user model.User

	phone := c.Query("phone")
	bookmark_name := c.Query("bookmark_name")
	if err := db.MysqlDB.First(&user, "phone = ?", phone).Error; err != nil {
		//返回结果
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败,根据phone参数找不到用户"})
		return
	}

	if err := db.MysqlDB.Model(&user).Update("bookmark_name", bookmark_name).Error; err != nil {
		//返回结果
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败"})
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
	if err := db.MysqlDB.Where("id = ?", docid).First(&docs).Error; err != nil {
		resp := gin.H{
			"message": "添加失败",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	caption := docs.Caption
	newBookmark := model.Bookmark{
		Phone:   phone,
		DocId:   docid,
		Caption: caption,
	}

	if err := db.MysqlDB.Create(&newBookmark).Error; err != nil {
		resp := gin.H{
			"message": "添加书签失败",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{"message": "OK"})

}

//删除单个书签
func Delete_bookmark(c *gin.Context) {

	var bookmark model.Bookmark

	phone := c.Query("phone")
	docid := c.Query("docid")

	if err := db.MysqlDB.Unscoped().Where("phone = ? and doc_id = ? ", phone, docid).Delete(&bookmark).Error; err != nil {
		resp := gin.H{
			"message": "删除书签失败",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	//返回结果
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})

}

//删除所有书签
func DeleteAll_bookmark(c *gin.Context) {

	var bookmark model.Bookmark
	var user model.User

	phone := c.Query("phone")

	if err := db.MysqlDB.First(&user, "phone = ?", phone).Error; err != nil {
		resp := gin.H{
			"message": "删除书签失败,根据phone参数找不到指定用户",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	if err := db.MysqlDB.Model(&user).Update("bookmark_name", "").Error; err != nil {
		resp := gin.H{
			"message": "用户书签名置空失败",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	if err := db.MysqlDB.Unscoped().Where("phone = ? ", phone).Delete(&bookmark).Error; err != nil {
		resp := gin.H{
			"message": "用户书书签删除失败",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	//返回结果
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})

}

//获得用户收藏的书签内容
func Get_bookmark(c *gin.Context) {

	phone := c.Query("phone")

	var bookmarks = []model.Bookmark{}

	if err := db.MysqlDB.Where("phone = ?", phone).Find(&bookmarks).Error; err != nil {
		resp := gin.H{
			"message": "用户书书签获取失败",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	// 写入Session
	setCurrentBookmark(c, bookmarks)

	c.JSON(http.StatusOK, gin.H{"bookmarks": bookmarks})

}
