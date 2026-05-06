package services

import (
	"fmt"
)

func SendWelcomeEmail(to, password, role string) error {
	subject := "Welcome to UBAPS - Your Account Details"
	
	body := fmt.Sprintf(`
		<div style="font-family: sans-serif; line-height: 1.6; color: #1e293b; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e2e8f0; border-radius: 12px;">
			<h2 style="color: #2563eb; margin-top: 0;">Welcome to UBAPS!</h2>
			<p>Hello,</p>
			<p>An administrator has created an account for you on the <strong>Bursary Award Processing System (UBAPS)</strong>. You can now log in using the credentials below:</p>
			
			<div style="background-color: #f8fafc; padding: 15px; border-radius: 8px; border: 1px solid #e2e8f0; margin: 20px 0;">
				<p style="margin: 5px 0;"><strong>Email:</strong> %s</p>
				<p style="margin: 5px 0;"><strong>Initial Password:</strong> %s</p>
				<p style="margin: 5px 0;"><strong>Assigned Role:</strong> %s</p>
			</div>
			
			<p>For security reasons, we recommend that you change your password immediately after your first login.</p>
			
			<div style="text-align: center; margin: 30px 0;">
				<a href="http://localhost:8080/login" style="display: inline-block; padding: 12px 24px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 6px; font-weight: 600;">Login to Your Account</a>
			</div>
			
			<hr style="border: 0; border-top: 1px solid #e2e8f0; margin: 20px 0;">
			<p style="font-size: 12px; color: #64748b;">This is an automated message. Please do not reply to this email.</p>
		</div>
	`, to, password, role)

	return SendEmail("richardsambo94@gmail.com", subject, body)
}
