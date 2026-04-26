package notifications

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// User_Created inserts a notification for a newly created user.
// It uses the provided transaction tx.
func Send_notification(user_id int64, tx pgx.Tx,message string,title string) error {
	query := `
	INSERT INTO notifications (
		recipient_user_id,
		notification_type,
		message,
		title
	)
	VALUES($1, $2, $3, $4)
	`

	_, err := tx.Exec(
		context.Background(),
		query,
		user_id,
		"SYSTEM",
		message,
		title,
	)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)
	}
	fmt.Println("success4")
	return nil
}
