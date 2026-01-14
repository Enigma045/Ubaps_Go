package user_logs

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func Create_user_log(tx pgx.Tx,
	userID *int64, // pointer allows NULL for system actions
	userRole string, // role at the time of action
	action string, // what action was performed
	target string, // affected entity, e.g., "user:42"
	status string, // e.g., "SUCCESS" or "FAILED"
	duration time.Duration, // how long the action took
) error {
	query := `
	INSERT INTO audit_user_logs (
		occurred_at,
		user_id,
		user_role,
		action,
		target,
		status,
		duration_ms
		)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 `
	_, err := tx.Exec(
		context.Background(),
		query,
		time.Now().UTC(),
		userID,
		userRole,
		action,
		target,
		status,
		int(duration.Milliseconds()),
	)

	return err
}
