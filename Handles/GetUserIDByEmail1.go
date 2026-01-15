package Handles

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// GetUserIDByEmail retrieves the user_id for a given email using the provided transaction.
// Returns an error if the user is not found or if the query fails.
func GetUserIDByEmail(email string, tx pgx.Tx) (int64, error) {
	if email == "" {
		return 0, fmt.Errorf("email cannot be empty")
	}

	var userID int64
	query := `SELECT user_id FROM users WHERE email = $1 LIMIT 1`

	err := tx.QueryRow(context.Background(), query, email).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("user not found for email: %s", email)
		}
		return 0, fmt.Errorf("failed to get user ID: %w", err)
	}
	fmt.Println("success3")
	return userID, nil
}

func GetEmailByUserID(userID int64, tx pgx.Tx) (string, error) {
	if userID <= 0 {
		return "", fmt.Errorf("userID must be a positive integer")
	}

	var email string
	query := `SELECT email FROM users WHERE user_id = $1 LIMIT 1`

	err := tx.QueryRow(context.Background(), query, userID).Scan(&email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("user not found for UserId: %d", userID)
		}
		return "", fmt.Errorf("failed to get user ID: %w", err)
	}
	fmt.Println("success3")
	return email, nil
}

func GetRegNumberFromEmail(email string) string {
	if email == "" {
		return ""
	}

	// Find the part before '@'
	parts := strings.Split(email, "@")
	if len(parts) == 0 {
		return ""
	}

	return parts[0]
}

// func GetUserIDFromContext(r *http.Request) (int64, error) {
// 	val := r.Context().Value(auth.userIDKey)
// 	if val == nil {
// 		return 0, fmt.Errorf("user not authenticated")
// 	}

// 	userID, ok := val.(int64)
// 	if !ok {
// 		return 0, fmt.Errorf("invalid userID type")
// 	}

// 	return userID, nil
// }

// Convert string to *string
func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Convert string to *float64
func FloatPtr(s string) (*float64, error) {
	if s == "" {
		return nil, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// Convert string to *time.Time (for date fields)
func TimePtr(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s) // standard HTML date format
	if err != nil {
		return nil, err
	}
	return &t, nil
}
