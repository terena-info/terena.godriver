package gomail

import (
	"fmt"
	"log"
	"net/smtp"
)

var (
	Gomail smtp.Auth
)

// GOMAIL
type GomailService struct {
	Mail   string
	Secret string
}

type Option struct {
	Html string
	To   string
}

func (opts GomailService) Send(mailOpts *Option) {
	to := []string{mailOpts.To}
	smtpPort := "587"
	smtpHost := "smtp.gmail.com"
	from := opts.Mail
	password := opts.Secret

	// Create authentication
	mail := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: Terena Authentication Agent\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
		"%s\r\n", opts.Mail, mailOpts, mailOpts.Html))

	err := smtp.SendMail(smtpHost+":"+smtpPort, mail, opts.Mail, to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
