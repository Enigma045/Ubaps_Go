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
	"strings"
)

type LogFilters struct {
	UserID       *int64
	TargetUserID *int64
	StartDate    *time.Time
	EndDate      *time.Time
}

func buildFilterQuery(baseQuery string, filters LogFilters) (string, []interface{}) {
	query := baseQuery
	params := []interface{}{}
	paramIdx := 1

	if filters.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", paramIdx)
		params = append(params, *filters.UserID)
		paramIdx++
	}
	if filters.TargetUserID != nil {
		query += fmt.Sprintf(" AND target_user_id = $%d", paramIdx)
		params = append(params, *filters.TargetUserID)
		paramIdx++
	}
	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND occurred_at >= $%d", paramIdx)
		params = append(params, *filters.StartDate)
		paramIdx++
	}
	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND occurred_at <= $%d", paramIdx)
		params = append(params, *filters.EndDate)
		paramIdx++
	}

	return query, params
}

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

type UnifiedLog struct {
	OccurredAt   utils.AutoTime  `json:"occurred_at"`
	UserID       *int64          `json:"user_id"`
	UserRole     string          `json:"user_role"`
	Action       string          `json:"action"`
	Target       string          `json:"target"`
	Status       string          `json:"status"`
	Duration     int64           `json:"duration_ms"`
	TargetUserID *int64          `json:"target_user_id"`
	Application  sql.NullString  `json:"application"`
	Amount       sql.NullFloat64 `json:"amount"`
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
	limit, offset int,
	filters LogFilters) ([]UserLog, int, error) {

	countBase := `SELECT COUNT(*) FROM audit_user_logs WHERE application IS NULL AND amount IS NULL`
	countQuery, countParams := buildFilterQuery(countBase, filters)
	
	var total int
	err := pool.QueryRow(ctx, countQuery, countParams...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	queryBase := `
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
`
	query, params := buildFilterQuery(queryBase, filters)
	query += fmt.Sprintf(" ORDER BY occurred_at DESC LIMIT $%d OFFSET $%d", len(params)+1, len(params)+2)
	params = append(params, limit, offset)

	rows, err := pool.Query(ctx, query, params...)
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
	limit, offset int,
	filters LogFilters) ([]PaymentLog, int, error) {

	countBase := `SELECT COUNT(*) FROM audit_user_logs WHERE amount IS NOT NULL`
	countQuery, countParams := buildFilterQuery(countBase, filters)
	
	var total int
	err := pool.QueryRow(ctx, countQuery, countParams...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	queryBase := `
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
`
	query, params := buildFilterQuery(queryBase, filters)
	query += fmt.Sprintf(" ORDER BY occurred_at DESC LIMIT $%d OFFSET $%d", len(params)+1, len(params)+2)
	params = append(params, limit, offset)

	rows, err := pool.Query(ctx, query, params...)
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
	limit, offset int,
	filters LogFilters) ([]ApplicationLog, int, error) {

	countBase := `SELECT COUNT(*) FROM audit_user_logs WHERE application IS NOT NULL`
	countQuery, countParams := buildFilterQuery(countBase, filters)
	
	var total int
	err := pool.QueryRow(ctx, countQuery, countParams...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	queryBase := `
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
`
	query, params := buildFilterQuery(queryBase, filters)
	query += fmt.Sprintf(" ORDER BY occurred_at DESC LIMIT $%d OFFSET $%d", len(params)+1, len(params)+2)
	params = append(params, limit, offset)

	rows, err := pool.Query(ctx, query, params...)
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

func Get_Unified_Logs(
	pool *pgxpool.Pool,
	ctx context.Context,
	limit, offset int,
	types []string,
	filters LogFilters) ([]UnifiedLog, int, error) {

	typeCond := ""
	if len(types) > 0 {
		conds := []string{}
		for _, t := range types {
			switch t {
			case "applications":
				conds = append(conds, "application IS NOT NULL")
			case "payments":
				conds = append(conds, "amount IS NOT NULL")
			case "users":
				conds = append(conds, "(application IS NULL AND amount IS NULL)")
			}
		}
		if len(conds) > 0 {
			typeCond = " AND (" + strings.Join(conds, " OR ") + ")"
		}
	}

	countBase := `SELECT COUNT(*) FROM audit_user_logs WHERE 1=1` + typeCond
	countQuery, countParams := buildFilterQuery(countBase, filters)

	var total int
	err := pool.QueryRow(ctx, countQuery, countParams...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	queryBase := `
    SELECT 
        occurred_at,
        user_id,
        user_role,
        action,
        target,
        status,
        duration_ms,
		target_user_id,
		application,
		amount
    FROM audit_user_logs
    WHERE 1=1
` + typeCond

	query, params := buildFilterQuery(queryBase, filters)
	query += fmt.Sprintf(" ORDER BY occurred_at DESC LIMIT $%d OFFSET $%d", len(params)+1, len(params)+2)
	params = append(params, limit, offset)

	rows, err := pool.Query(ctx, query, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []UnifiedLog
	for rows.Next() {
		var log UnifiedLog
		err := rows.Scan(
			&log.OccurredAt,
			&log.UserID,
			&log.UserRole,
			&log.Action,
			&log.Target,
			&log.Status,
			&log.Duration,
			&log.TargetUserID,
			&log.Application,
			&log.Amount,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, nil
}