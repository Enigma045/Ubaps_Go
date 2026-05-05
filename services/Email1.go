package services

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
)

func SendEmail(to, subject, body string) error {
	from := "richardsambo94@gmail.com"
	password := "opsi gofk ezyy pzal"

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

func SendEmailWithAttachment(to, subject, body, fileName string, fileContent []byte) error {
	from := "richardsambo94@gmail.com"
	password := "opsi gofk ezyy pzal"

	boundary := "my-boundary-12345"

	header := fmt.Sprintf("From: %s\r\n", from)
	header += fmt.Sprintf("To: %s\r\n", to)
	header += fmt.Sprintf("Subject: %s\r\n", subject)
	header += "MIME-Version: 1.0\r\n"
	header += fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary)
	header += "\r\n"

	// Body part
	message := fmt.Sprintf("--%s\r\n", boundary)
	message += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	message += "\r\n"
	message += body + "\r\n"
	message += "\r\n"

	// Attachment part
	message += fmt.Sprintf("--%s\r\n", boundary)
	message += fmt.Sprintf("Content-Type: application/octet-stream; name=\"%s\"\r\n", fileName)
	message += "Content-Transfer-Encoding: base64\r\n"
	message += fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName)
	message += "\r\n"

	// Encode content to base64
	b := make([]byte, base64.StdEncoding.EncodedLen(len(fileContent)))
	base64.StdEncoding.Encode(b, fileContent)

	message += string(b) + "\r\n"
	message += fmt.Sprintf("--%s--", boundary)

	return smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from,
		[]string{to},
		[]byte(header+message),
	)
}
