package services

import (
	"fmt"
	"net/smtp"
)

func SendEmail(to, subject, body string) error {
	from := "cen-01-14-22@unilia.ac.mw"
	password := "rrfa fthh yngg lcad"

	msg := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body
	fmt.Println("success1")
	return smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from,
		[]string{to},
		[]byte(msg),
	)
}
