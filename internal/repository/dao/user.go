package dao

import (
	"context"
	"gorm.io/gorm"
)

type User struct {
	Id       int64 `gorm:"primaryKey"`
	Email    string
	Password string
	Ctime    string
	Utime    string
}

type UserDao struct {
	db *gorm.DB
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	return dao.db.WithContext(ctx).Create(&u).Error
}
