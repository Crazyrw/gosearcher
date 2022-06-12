package model

type Docs struct {
	ID      int    `gorm:"type:int;primaryKey;autoIncrement"`
	Url     string `gorm:"type:varchar(500)"`
	Caption string `gorm:"type:varchar(500)"`
}
