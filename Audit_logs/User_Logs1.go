package user_logs

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"ubaps/Db"
	"ubaps/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserLog struct {
	OccurredAt   utils.AutoTime `json:"occurred_at"`
	UserID       *int64         `json:"user_id"` // performer
	UserRole     string         `json:"user_role"`
	Action       string         `json:"action"`
	Target       string         `json:"target"`
	Status       string         `json:"status"`
	Duration     int64          `json:"duration_ms"`
	TargetUserID *int64         `json:"target_user_id"`
}

type PaymentLog struct {
	OccurredAt   utils.AutoTime `json:"occurred_at"`
	UserID       *int64         `json:"user_id"` // performer
	UserRole     string         `json:"user_role"`
	Action       string         `json:"action"`
	Target       string         `json:"target"`
	Status       string         `json:"status"`
	Duration     int64          `json:"duration_ms"`
	Application  sql.NullString `json:"application"`
	Amount       float64        `json:"amount"`
	TargetUserID *int64         `json:"target_user_id"`
}

type ApplicationLog struct {
	OccurredAt   utils.AutoTime  `json:"occurred_at"`
	UserID       *int64          `json:"user_id"` // performer
	UserRole     string          `json:"user_role"`
	Action       string          `json:"action"`
	Target       string          `json:"target"`
	Status       string          `json:"status"`
	Duration     int64           `json:"duration_ms"`
	Application  string          `json:"application"`
	Amount       sql.NullFloat64 `json:"amount"`
	TargetUserID *int64          `json:"target_user_id"`
}

func Create_user_log(tx pgx.Tx,
	userID *int64, // pointer allows NULL for system actions
	userRole string, // role at the time of action
	action string, // what action was performed
	target string, // affected entity, e.g., "user:42"
	status string, // e.g., "SUCCESS" or "FAILED"
	duration time.Duration, // how long the action took
	targetUserID *int64,
) error {
	query := `
	INSERT INTO audit_user_logs (
		occurred_at,
		user_id,
		user_role,
		action,
		target,
		status,
		duration_ms,
		target_user_id
		)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 `
	
	now := time.Now().UTC()
	dur := int(duration.Milliseconds())

	if tx != nil {
		_, err := tx.Exec(context.Background(), query, now, userID, userRole, action, target, status, dur, targetUserID)
		return err
	}
	
	_, err := Db.DB.Exec(context.Background(), query, now, userID, userRole, action, target, status, dur, targetUserID)
	return err
}

func Create_application_log(tx pgx.Tx,
	userID *int64,
	userRole string,
	action string,
	target string,
	status string,
	duration time.Duration,
	applicationName string,
	amount *float64,
	targetUserID *int64,
) error {
	query := `
	INSERT INTO audit_user_logs (
		occurred_at,
		user_id,
		user_role,
		action,
		target,
		status,
		duration_ms,
		application,
		amount,
		target_user_id
		)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 `
	
	now := time.Now().UTC()
	dur := int(duration.Milliseconds())

	pID := "SYSTEM"
	if userID != nil {
		pID = fmt.Sprintf("%d", *userID)
	}
	tID := "SYSTEM"
	if targetUserID != nil {
		tID = fmt.Sprintf("%d", *targetUserID)
	}
	target = fmt.Sprintf("USER:%s->APP:%s", pID, tID)

	if tx != nil {
		_, err := tx.Exec(context.Background(), query, now, userID, userRole, action, target, status, dur, applicationName, amount, targetUserID)
		return err
	}

	_, err := Db.DB.Exec(context.Background(), query, now, userID, userRole, action, target, status, dur, applicationName, amount, targetUserID)
	return err
}

func Create_payment_log(tx pgx.Tx,
	userID *int64,
	userRole string,
	action string,
	target string,
	status string,
	duration time.Duration,
	amount float64,
	application string,
	targetUserID *int64,
) error {
	query := `
	INSERT INTO audit_user_logs (
		occurred_at,
		user_id,
		user_role,
		action,
		target,
		status,
		duration_ms,
		amount,
		application,
		target_user_id
		)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 `
	
	now := time.Now().UTC()
	dur := int(duration.Milliseconds())

	pID := "SYSTEM"
	if userID != nil {
		pID = fmt.Sprintf("%d", *userID)
	}
	tID := "SYSTEM"
	if targetUserID != nil {
		tID = fmt.Sprintf("%d", *targetUserID)
	}
	target = fmt.Sprintf("USER:%s->APP:%s", pID, tID)

	if tx != nil {
		_, err := tx.Exec(context.Background(), query, now, userID, userRole, action, target, status, dur, amount, application, targetUserID)
		return err
	}

	_, err := Db.DB.Exec(context.Background(), query, now, userID, userRole, action, target, status, dur, amount, application, targetUserID)
	return err
}

func Get_User_Logs(
	pool *pgxpool.Pool,
	ctx context.Context,
	limit, offset int) ([]UserLog, int, error) {

	countQuery := `SELECT COUNT(*) FROM audit_user_logs WHERE application IS NULL AND amount IS NULL`
	var total int
	err := pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
    SELECT 
        occurred_at,
        user_id,
        user_role,
        action,
        target,
        status,
        duration_ms,
		target_user_id
    FROM audit_user_logs
    WHERE application IS NULL
      AND amount IS NULL
    ORDER BY occurred_at DESC
    LIMIT $1 OFFSET $2
`
	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []UserLog
	for rows.Next() {
		var log UserLog
		err := rows.Scan(
			&log.OccurredAt,
			&log.UserID,
			&log.UserRole,
			&log.Action,
			&log.Target,
			&log.Status,
			&log.Duration,
			&log.TargetUserID,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, nil
}

func Get_Payment_Logs(
	pool *pgxpool.Pool,
	ctx context.Context,
	limit, offset int) ([]PaymentLog, int, error) {

	countQuery := `SELECT COUNT(*) FROM audit_user_logs WHERE amount IS NOT NULL`
	var total int
	err := pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
    SELECT 
        occurred_at,
        user_id,
        user_role,
        action,
        target,
        status,
        duration_ms,
		application,
		amount,
		target_user_id
    FROM audit_user_logs
    WHERE amount IS NOT NULL
    ORDER BY occurred_at DESC
    LIMIT $1 OFFSET $2
`
	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []PaymentLog
	for rows.Next() {
		var log PaymentLog
		err := rows.Scan(
			&log.OccurredAt,
			&log.UserID,
			&log.UserRole,
			&log.Action,
			&log.Target,
			&log.Status,
			&log.Duration,
			&log.Application,
			&log.Amount,
			&log.TargetUserID,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, nil
}

func Get_Application_Logs(
	pool *pgxpool.Pool,
	ctx context.Context,
	limit, offset int) ([]ApplicationLog, int, error) {

	countQuery := `SELECT COUNT(*) FROM audit_user_logs WHERE application IS NOT NULL`
	var total int
	err := pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
    SELECT 
        occurred_at,
        user_id,
        user_role,
        action,
        target,
        status,
        duration_ms,
		application,
		amount,
		target_user_id
    FROM audit_user_logs
    WHERE application IS NOT NULL
    ORDER BY occurred_at DESC
    LIMIT $1 OFFSET $2
`
	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []ApplicationLog
	for rows.Next() {
		var log ApplicationLog
		err := rows.Scan(
			&log.OccurredAt,
			&log.UserID,
			&log.UserRole,
			&log.Action,
			&log.Target,
			&log.Status,
			&log.Duration,
			&log.Application,
			&log.Amount,
			&log.TargetUserID,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, nil
}