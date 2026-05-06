package services

import (
	"fmt"
)

func SendPasswordResetEmail(to, token string) error {
	link := fmt.Sprintf(
		"http://localhost:8080/reset-password?token=%s",
		token,
	)

	body := fmt.Sprintf(`
		<h2>Reset Your UBAPS Password</h2>
		<p>You requested a password reset. Click the link below to set a new password:</p>
		<a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 5px;">Reset Password</a>
		<p>This link will expire in 1 hour.</p>
		<p>If you did not request this, please ignore this email.</p>
	`, link)

	return SendEmail(to, "Password Reset Request", body)
}
