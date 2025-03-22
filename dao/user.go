package dao

import (
	"errors"
	"go-im/db"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

var dbIns *gorm.DB

type User struct {
	Id       uint `gorm:"primary_key"`
	UserName string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func Init(){
	dbIns = db.InitDB()
}

func (u *User) Add() (userId uint, err error) {
	if u.UserName == "" || u.Password == "" {
		return 0, errors.New("user_name or password empty!")
	}
	oUser := u.CheckHaveUserName(u.UserName)
	if oUser.Id > 0 {
		return oUser.Id, nil
	}
	if err = dbIns.Create(&u).Error; err != nil {
		return 0, err
	}

	return u.Id, nil
}

func (u *User) CheckHaveUserName(userName string) (data User) {
	if err := dbIns.Where("username = ?", userName).Take(&data); err != nil {
		logx.Infof("404 not found")
		return
	}
	return
}