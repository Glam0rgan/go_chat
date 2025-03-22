package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (dbIns *gorm.DB){
	dsn := "root:Bzt1453529.@tcp(127.0.0.1:3306)/go_chat"
	dbIns, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("dsn:%s invalid,err:%v\n", dsn, err)
		return
	}

	fmt.Println("connect ok")

	return dbIns
}