package main

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
)

const (
	host     = "smtp.yandex.ru"
	port     = 465
	from     = "noreply@tvbit.co"
	to       = "goo007@mail.ru"
	password = "xSVMr7WHpszf7DWs"
	subject  = "test subject 4"
	body     = "test body 4"
)

func main() {
	from1 := mail.Address{Address: from}
	to := mail.Address{Address: to}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from1.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := fmt.Sprintf("%s:%d", host, port)

	auth := smtp.PlainAuth("", from, password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		panic(err)
	}

	// To && From
	if err = c.Mail(from1.Address); err != nil {
		panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		panic(err)
	}

	err = w.Close()
	if err != nil {
		panic(err)
	}

	if err := c.Quit(); err != nil {
		panic(err)
	}
}
