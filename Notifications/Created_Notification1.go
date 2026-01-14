package notifications

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// User_Created inserts a notification for a newly created user.
// It uses the provided transaction tx.
func User_Created(user_id int64, tx pgx.Tx) error {
	query := `
	INSERT INTO notifications (
		recipient_user_id,
		notification_type,
		message
	)
	VALUES($1, $2, $3)
	`

	_, err := tx.Exec(
		context.Background(),
		query,
		user_id,
		"SYSTEM",
		"Your account has been created.",
	)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)
	}
	fmt.Println("success4")
	return nil
}
