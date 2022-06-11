package controller

import (
	"crypto/md5"
	"fmt"
	"goSearcher/searcher/db"
	"goSearcher/web/model"
	"net/http"
	"regexp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//func DbConnect() (db *gorm.DB) {
//	dsn := "ligen:LiGen1129!@tcp(127.0.0.1:3306)/goSearcher?charset=utf8mb4&parseTime=True&loc=Local"
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	return db
//}

func setCurrentUser(c *gin.Context, userInfo model.User) {
	session := sessions.Default(c)
	session.Set("currentUser", userInfo)
	// 一定要Save否则不生效，若未使用gob注册User结构体，调用Save时会返回一个Error
	session.Save()
}

func getCurrentUser(c *gin.Context) (userInfo model.User) {
	session := sessions.Default(c)
	userInfo = session.Get("currentUser").(model.User) // 类型转换一下
	return userInfo
}

func delCurrentUser(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("currentUser")
	session.Save()
}

func UserLoginGet(c *gin.Context) {
	resp := gin.H{"message": "OK"}
	c.HTML(http.StatusOK, "login.tmpl", resp)
}

func UserLoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 返回结果
	resp := gin.H{}

	data := []byte(password)
	has := md5.Sum(data)
	new_password := fmt.Sprintf("%x", has)

	userObject := model.User{}

	// 查询数据库
	if err := db.MysqlDB.Where("username = ? AND password = ?", username, new_password).First(&userObject).Error; err != nil {
		resp = gin.H{
			"login":   false,
			"message": "用户名或密码错误",
		}
		c.HTML(http.StatusBadRequest, "login.tmpl", resp)
		return
	}

	// 写入Session
	setCurrentUser(c, userObject)
	userInfo := getCurrentUser(c)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"userInfo": userInfo})

}

func UserRegisterGet(c *gin.Context) {
	resp := gin.H{"message": "OK"}
	c.HTML(http.StatusOK, "register.tmpl", resp)
}

func UserRegisterPost(c *gin.Context) {
	username := c.PostForm("username")
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	confirm_password := c.PostForm("confirm_password")

	// 返回结果
	resp := gin.H{"message": "OK"}

	// 查询用户名是否存在
	if err := db.MysqlDB.Where("username = ?", username).First(&model.User{}).Error; err == nil {
		resp = gin.H{
			"message": "用户名已存在",
		}
		c.HTML(http.StatusBadRequest, "register.tmpl", resp)

		return
	}

	// 判断手机号是否规范
	matched, err := regexp.MatchString("^1[3456789]\\d{9}$", phone)
	if matched != true || err != nil {
		resp = gin.H{
			"message": "手机号不符合规范",
		}
		c.HTML(http.StatusBadRequest, "register.tmpl", resp)

		return
	}

	// 查询手机号是否存在
	if err := db.MysqlDB.Where("phone = ?", phone).First(&model.User{}).Error; err == nil {
		resp = gin.H{
			"message": "手机号已存在",
		}
		c.HTML(http.StatusBadRequest, "register.tmpl", resp)

		return
	}

	// 重复输入密码错误
	if password != confirm_password {
		resp = gin.H{
			"message": "两次密码不一致",
		}
		c.HTML(http.StatusBadRequest, "register.tmpl", resp)
		return
	}

	// 密码加密
	data := []byte(password)
	has := md5.Sum(data)
	new_password := fmt.Sprintf("%x", has) //将[]byte转成16进制

	user := model.User{Username: username, Phone: phone, Password: new_password}

	result := db.MysqlDB.Create(&user) // 通过数据的指针来创建
	fmt.Println(result)

	c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"message": "OK"})

}

func UserLogout(c *gin.Context) {
	delCurrentUser(c)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

//用户注销账号
func UserDelete(c *gin.Context) {
	var user model.User
	var bookmark model.Bookmark

	phone := c.Query("phone")

	db.MysqlDB.Unscoped().Where("telephone = ?", phone).Delete(&user)
	//用户注销的时候，连带把他所有的书签信息全部删除
	db.MysqlDB.Unscoped().Where("telephone = ? ", phone).Delete(&bookmark)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}
