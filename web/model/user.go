package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
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
