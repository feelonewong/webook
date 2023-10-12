package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB, tables interface{}) error {
	return db.AutoMigrate(tables)
}
