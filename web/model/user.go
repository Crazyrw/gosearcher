package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       int    `gorm:"primary_key"`
	Username string `gorm:"not null"`
	Phone    string `gorm:"not null"`
	Password string `gorm:"not null"`
}

type UserSave struct {
	User      User `gorm:"ForeignKey:UserId;AssociationForeignKey:ID"`
	UserId    int
	SaveDocId int `gorm:"not null"`
}

<<<<<<< HEAD
func UserDBInit() {
	dsn := "root:122513gzhGZH!!@tcp(decs.pcl.ac.cn:1762)/search_engine?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserSave{})

	fmt.Println("Finish migrate.")

}
=======
//func UserDBInit() {
//	dsn := "ligen:LiGen1129!@tcp(127.0.0.1:3306)/goSearcher?charset=utf8mb4&parseTime=True&loc=Local"
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	db.AutoMigrate(&User{})
//	db.AutoMigrate(&UserSave{})
//
//	fmt.Println("Finish migrate.")
//
//}
>>>>>>> 39830c156187479b44a4133abe0e0253d3952e65
