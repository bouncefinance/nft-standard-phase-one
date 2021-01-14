package models

import (
	"Ankr-gin-ERC721/pkg/logger"
	"Ankr-gin-ERC721/pkg/setting"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var (
		err                          error
		dbName, user, password, host string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		panic( err)
	}

	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	password = password + "#"

	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName)
	db, err = gorm.Open(mysql.Open(str), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("mysql link", str).Msg("")

	d, err := db.DB()
	if err != nil {
		panic(err)
	}
	d.SetMaxIdleConns(10)
	d.SetMaxOpenConns(100)
}

func CloseDB() {
	//defer db.Close()
}
