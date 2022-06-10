package model

type Query struct {
	ID     int    `gorm:"primary_key"`
	Query  string `gorm:"not null"`
	DocIds string `gorm:"not null"`
}
