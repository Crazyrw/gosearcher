package model

import "gorm.io/gorm"

type Bookmark struct {
	gorm.Model
	Phone   string `gorm:"not null"`
	DocId   string `gorm:"not null"`
	Caption string `gorm:"not null"`
}
