package utils

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Notification struct {
    Title string `json:"title"`
	Message string `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func ReciveNotifications(
	pool *pgxpool.Pool,
	 ctx context.Context,
	userID int64,
) ([]Notification, error) {

	query := `
		SELECT title, message, created_at
		FROM notifications
		WHERE recipient_user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var notify Notification

		err := rows.Scan(
			&notify.Title,
			&notify.Message,
			&notify.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notify)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}


func NotificationCounter(
	pool *pgxpool.Pool,
	ctx context.Context,
	userID int64,
) (int, error) {

	var count int

	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE recipient_user_id = $1
	`

	err := pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
