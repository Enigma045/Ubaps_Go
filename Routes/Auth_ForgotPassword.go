package Routes

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"ubaps/Db"
	"ubaps/services"
	user_logs "ubaps/Audit_logs"
	"ubaps/Notifications"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// EnsurePasswordResetsTable checks and creates the password_resets table if it doesn't exist.
func EnsurePasswordResetsTable(ctx context.Context) {
	query := `
	CREATE TABLE IF NOT EXISTS password_resets (
		id BIGSERIAL PRIMARY KEY,
		email TEXT NOT NULL REFERENCES users(email) ON DELETE CASCADE,
		token TEXT NOT NULL,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(email)
	);`
	_, err := Db.DB.Exec(ctx, query)
	if err != nil {
		log.Println("Error ensuring password_resets table:", err)
	}
}

func ForgotPassword_page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/public/ForgotPassword.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func ResetPassword_page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/public/ResetPassword.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func ForgotPasswordAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	EnsurePasswordResetsTable(ctx)

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	// Check if user exists and get details for logging
	var userID int64
	var role string
	err = tx.QueryRow(ctx, "SELECT user_id, user_type FROM users WHERE email = $1", email).Scan(&userID, &role)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			// Log failure (Email not found) - Use system as role since user doesn't exist
			user_logs.Create_user_log(tx, nil, "SYSTEM", "FORGOT_PASSWORD_FAILED", fmt.Sprintf("Reset attempt for non-existent email: %s", email), "FAILED", time.Since(start), nil)
			tx.Commit(ctx)
			
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"success": true, "message": "If that email exists, a reset link has been sent."}`)
			return
		}
		log.Println("Database error checking user:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	token := hex.EncodeToString(b)
	expiresAt := time.Now().Add(1 * time.Hour)

	// Save token
	query := `
	INSERT INTO password_resets (email, token, expires_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (email) DO UPDATE
	SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at;`
	
	_, err = tx.Exec(ctx, query, email, token, expiresAt)
	if err != nil {
		log.Println("Database error saving reset token:", err)
		user_logs.Create_user_log(tx, &userID, role, "FORGOT_PASSWORD_FAILED", "DB Error saving token", "FAILED", time.Since(start), &userID)
		tx.Commit(ctx)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Audit Log (Success)
	user_logs.Create_user_log(tx, &userID, role, "FORGOT_PASSWORD_REQUESTED", fmt.Sprintf("Password reset link generated for %s", email), "SUCCESS", time.Since(start), &userID)

	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send email (Hardcoded recipient as requested)
	// We do this after commit to ensure the token is saved
	err = services.SendPasswordResetEmail("richardsambo94@gmail.com", token)
	if err != nil {
		log.Println("Email error sending reset link:", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true, "message": "A reset link has been sent to your email."}`)
}

func ResetPasswordAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()
	token := r.FormValue("token")
	password := r.FormValue("password")

	if token == "" || password == "" {
		http.Error(w, "Token and password are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	
	// Start transaction for logging
	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	// Validate token and get user details for logging
	var email string
	var userID int64
	var role string
	var expiresAt time.Time
	
	err = tx.QueryRow(ctx, `
		SELECT pr.email, u.user_id, u.user_type, pr.expires_at 
		FROM password_resets pr
		JOIN users u ON u.email = pr.email
		WHERE pr.token = $1
	`, token).Scan(&email, &userID, &role, &expiresAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			user_logs.Create_user_log(tx, nil, "SYSTEM", "PASSWORD_RESET_FAILED", "Invalid or missing token", "FAILED", time.Since(start), nil)
			tx.Commit(ctx)
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
		}
		log.Println("Database error validating token:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if time.Now().After(expiresAt) {
		user_logs.Create_user_log(tx, &userID, role, "PASSWORD_RESET_FAILED", "Token expired", "FAILED", time.Since(start), &userID)
		tx.Commit(ctx)
		http.Error(w, "Token has expired", http.StatusBadRequest)
		return
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update user password
	_, err = tx.Exec(ctx, "UPDATE users SET password_hash = $1 WHERE email = $2", string(hash), email)
	if err != nil {
		user_logs.Create_user_log(tx, &userID, role, "PASSWORD_RESET_FAILED", "DB Error updating password", "FAILED", time.Since(start), &userID)
		tx.Commit(ctx)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Delete reset token
	_, err = tx.Exec(ctx, "DELETE FROM password_resets WHERE token = $1", token)
	if err != nil {
		log.Println("Error deleting token:", err)
		// Not fatal, but good to know
	}

	// Audit Log (Success)
	user_logs.Create_user_log(tx, &userID, role, "PASSWORD_RESET_SUCCESS", "User successfully reset their password via email link", "SUCCESS", time.Since(start), &userID)

	// In-app Notification
	notifications.Send_notification(userID, tx, "Your password has been successfully reset via email link.", "Password Reset Successful")

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true, "message": "Password updated successfully."}`)
}
