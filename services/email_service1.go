package services

import "fmt"

func SendVerificationEmail(to, token string) error {
	link := fmt.Sprintf(
		"http://localhost:8080/verify-email?token=%s",
		token,
	)

	body := fmt.Sprintf(`
		<h2>Verify your UBAPS account</h2>
		<p>Click the link below to verify your email:</p>
		<a href="%s">Verify Email</a>
	`, link)
	fmt.Println("success1")
	return SendEmail(to, "Verify your email", body)
}
