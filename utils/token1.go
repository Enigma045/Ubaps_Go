package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// GenerateVerificationToken creates a unique token for the given email.
// It inserts or updates the token in the email_verifications table within the given transaction.
func GenerateVerificationToken(email string, tx pgx.Tx) (string, error) {
	// 1️⃣ Check email
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}

	if !strings.Contains(email, "@") {
		return "", fmt.Errorf("invalid email format")
	}

	// 2️⃣ Check transaction
	if tx == nil {
		return "", fmt.Errorf("transaction is nil")
	}

	log.Println("✅ TOKEN STEP 1: input validated")

	// 3️⃣ Generate token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("token generation failed: %w", err)
	}

	token := hex.EncodeToString(bytes)
	expiresAt := time.Now().Add(24 * time.Hour)

	log.Println("✅ TOKEN STEP 2: token generated")
	log.Println(email)
	log.Println(token)
	log.Println(expiresAt)
	// 4️⃣ Insert / update
	query := `
	INSERT INTO email_verifications (email, token, expires_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (email)
	DO UPDATE
	SET token = EXCLUDED.token,
	    expires_at = EXCLUDED.expires_at
	`

	ct, err := tx.Exec(
		context.Background(),
		query,
		email,
		token,
		expiresAt,
	)

	if err != nil {
		return "", fmt.Errorf("email_verifications insertion failed: %w", err)
	}
	fmt.Println("Rows affected:", ct.RowsAffected())
	// 5️⃣ Check rows affected
	if ct.RowsAffected() == 0 {
		return "", fmt.Errorf("no rows affected when inserting verification token")
	}

	log.Println("✅ TOKEN STEP 3: database updated")

	return token, nil
}
