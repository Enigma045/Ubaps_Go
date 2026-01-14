package Handles

import (
	"context"
	"fmt"

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
