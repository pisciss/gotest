package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func Connect() {
	d, err := gorm.Open("mysql", "root:@(localhost:3306)/golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db = d

}
func GetDB() *gorm.DB {
	return db
}

func CheckConnection() int {

	pingErr := db.DB().Ping()
	if pingErr != nil {
		return 0
	} else {
		return 1
	}

}
