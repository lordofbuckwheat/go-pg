package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID uint64 `gorm:"PRIMARY_KEY"`
}

type User struct {
	Model
	Login     string `gorm:"NOT NULL"`
	Password  string `gorm:"NOT NULL"`
	AccountID uint64
	Account   *Account
}

type Account struct {
	Model
	Title string `gorm:"NOT NULL"`
	Users []User
}

func main() {
	db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/tvbit_test?charset=utf8mb4")
	if err != nil {
		panic("failed to connect database")
	}
	defer func() { _ = db.Close() }()
	db.LogMode(true)
	db.AutoMigrate(
		&Account{},
		&User{},
	)
	var account = &Account{
		Title: "account 1",
	}
	if err := db.Create(account).Error; err != nil {
		panic(err)
	}
	if err := db.Model(account).Association("Users").Replace([]User{{
		Login:    "user1",
		Password: "password1",
	}, {
		Login:    "user2",
		Password: "password2",
	}}).Error; err != nil {
		panic(err)
	}
	fmt.Println("account", MustMarshal(account))
	var user = &User{}
	if err := db.Find(user).Error; err != nil {
		panic(err)
	}
	fmt.Println("user", MustMarshal(user))
	var account2 = &Account{}
	if err := db.Model(user).Association("Account").Find(account2).Error; err != nil {
		panic(err)
	}
	fmt.Println("user", MustMarshal(user), MustMarshal(account2))
}

func MustMarshal(src interface{}) string {
	result, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	return string(result)
}
