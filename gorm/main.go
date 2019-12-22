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
	Profile   *Profile
}

type Profile struct {
	Model
	Name   string `gorm:"NOT NULL"`
	UserID uint64
	User   *User
}

type Account struct {
	Model
	Title string `gorm:"NOT NULL"`
	Users []*User
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
		&Profile{},
	)
	var account = &Account{
		Title: "account 1",
	}
	if err := db.Create(account).Error; err != nil {
		panic(err)
	}
	var users = []*User{{
		Login:    "user1",
		Password: "password1",
		Profile: &Profile{
			Name: "profile 1",
		},
		Account: account,
	}, {
		Login:    "user2",
		Password: "password2",
	}}
	if err := db.Model(account).Association("Users").Replace(users).Error; err != nil {
		panic(err)
	}
	fmt.Println("account", MustMarshal(account), MustMarshal(users))
	users[0].Profile = &Profile{
		Name: "profile 2",
	}
	if err := db.Save(users[0]).Error; err != nil {
		panic(err)
	}
	fmt.Println("user after replace", MustMarshal(users[0]))
	account.Users = []*User{{
		Login:    "user3",
		Password: "password3",
	}}
	if err := db.Save(account).Error; err != nil {
		panic(err)
	}
	fmt.Println("account", MustMarshal(account))
	account.Users = nil
	var user = &User{}
	if err := db.Find(user).Error; err != nil {
		panic(err)
	}
	fmt.Println("user", MustMarshal(user))
	var account2 = &Account{}
	if err := db.Model(user).Association("Account").Find(account2).Error; err != nil {
		panic(err)
	}
	user.Account = account2
	fmt.Println("user", MustMarshal(user), MustMarshal(account2))
	users = nil
	if err := db.Model(account2).Association("Users").Find(&users).Error; err != nil {
		panic(err)
	}
	fmt.Println("account2", MustMarshal(account2), MustMarshal(users))
}

func MustMarshal(src interface{}) string {
	result, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	return string(result)
}
