package user_logs

import (
	"context"
	"database/sql"
	"time"
	"ubaps/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserLog struct {
	OccurredAt utils.AutoTime `json:"occurred_at"`
	UserID     *int64 `json:"user_id"`// pointer allows NULL for system actions
    UserRole   string `json:"user_role"` // role at the time of action
    Action     string `json:"action"` // what action was performed
    Target    string `json:"target"` // affected entity, e.g., "user:42"
	Status    string `json:"status"` // e.g., "SUCCESS" or "FAILED"
	Duration  int64  `json:"duration_ms"` // duration in milliseconds
}

type PaymentLog struct {
	OccurredAt utils.AutoTime `json:"occurred_at"`
	UserID     *int64 `json:"user_id"`// pointer allows NULL for system actions
    UserRole   string `json:"user_role"` // role at the time of action
    Action     string `json:"action"` // what action was performed
    Target    string `json:"target"` // affected entity, e.g., "user:42"
	Status    string `json:"status"` // e.g., "SUCCESS" or "FAILED"
	Duration  int64  `json:"duration_ms"` // duration in milliseconds
	Application sql.NullString `json:"application"`
    Amount     float64 `json:"amount"`
}

type ApplicationLog struct {
	OccurredAt utils.AutoTime `json:"occurred_at"`
	UserID     *int64 `json:"user_id"`// pointer allows NULL for system actions
    UserRole   string `json:"user_role"` // role at the time of action
    Action     string `json:"action"` // what action was performed
    Target    string `json:"target"` // affected entity, e.g., "user:42"
	Status    string `json:"status"` // e.g., "SUCCESS" or "FAILED"
	Duration  int64  `json:"duration_ms"` // duration in milliseconds
	Application string `json:"application"`
    Amount     sql.NullFloat64 `json:"amount"`
}

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

func Get_User_Logs(
	 pool *pgxpool.Pool,
	 ctx context.Context) ([]UserLog, error){

	query := `
    SELECT 
        occurred_at,
        user_id,
        user_role,
        action,
        target,
        status,
        duration_ms
    FROM audit_user_logs
    WHERE application IS NULL
      AND amount IS NULL
`
 	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil,err
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
		)
		if err != nil {
			return nil,err
		}

		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return nil,err
	}

	return logs,nil
}

func Get_Payment_Logs(
	 pool *pgxpool.Pool,
	 ctx context.Context) ([]PaymentLog, error){

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
		amount
    FROM audit_user_logs
    WHERE amount IS NOT NULL
`
 	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil,err
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
		)
		if err != nil {
			return nil,err
		}

		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return nil,err
	}

	return logs,nil
}

func Get_Application_Logs(
	 pool *pgxpool.Pool,
	 ctx context.Context) ([]ApplicationLog, error){

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
		amount
    FROM audit_user_logs
    WHERE application IS NOT NULL
`
 	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil,err
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
		)
		if err != nil {
			return nil,err
		}

		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return nil,err
	}

	return logs,nil
}