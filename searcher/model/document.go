package model

// Docs 数据库存储的文档对象
type Docs struct {
	ID      int    `gorm:"type:int;primaryKey;autoIncrement"`
	Url     string `gorm:"type:varchar(500)"`
	Caption string `gorm:"type:varchar(500)"`
}
