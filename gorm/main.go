package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/tvbit_test?charset=utf8mb4")
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {_ = db.Close()}()
	fmt.Println("123")
}
