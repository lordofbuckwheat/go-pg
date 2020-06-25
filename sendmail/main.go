package main

import (
	"fmt"
	"net/smtp"
)

const (
	host     = "smtp.gmail.com"
	port     = 587
	from     = "lord.of.buckwheat@gmail.com"
	to       = "goo007@mail.ru"
	password = ""
	subject  = "test subject 2"
	body     = "test body 2"
)

func main() {
	auth := smtp.PlainAuth("", from, password, host)
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, from, []string{to}, []byte(msg))
	if err != nil {
		panic(err)
	}
	fmt.Println("done")
}
