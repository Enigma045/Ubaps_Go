package Handles

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type Requests struct {
	Id string `json:"id"`
	First  string `json:"first_name"`
	Last   string `json:"last_name"`
	Status string `json:"status"`
	Previous string `json:"previous_amount"`
	Amount string `json:"payment_amount"`
	Total_Requested_Amount string `json:"cumulative_amount"`
	Details sql.NullString `json:"details"`
	RegNum string `json:"registration_number"`
	ParentStatus string `json:"parent_guardian_status"`
	EmploymentStatus string `json:"employment_status"`
	Programme string `json:"programme"`
}

type StatementRequest struct {
	First            string `json:"first_name"`
	Last             string `json:"last_name"`
	RegNum           string `json:"registration_number"`
	ParentStatus     string `json:"parent_guardian_status"`
	EmploymentStatus string `json:"employment_status"`
	Income           string `json:"income"`
	Priority         string `json:"priority"`
	Status           string `json:"status"`
}

func GetRequest_Info(
	pool *pgxpool.Pool,
	ctx context.Context,
	request []string,
) ([]Requests, error) {

	// Safety: avoid SQL error on empty slice
	if len(request) == 0 {
		return []Requests{}, nil
	}


	query := 
	`
	SELECT
    u.name,
    u.surname ,
	a.request_id,
	a.previous_amount,
	a.payment_amount,
	a.cumulative_amount,
	a.details,
	a.request_status,
	COALESCE(app.registration_number, ''),
	COALESCE(app.parent_guardian_status, ''),
	COALESCE(app.guardian_employment_status, ''),
	COALESCE(app.programme, '')
    FROM financial_request a
    JOIN users u ON u.user_id = a.user_id
	LEFT JOIN applications app ON app.user_id = a.user_id
    WHERE a.request_status = ANY($1);
	`
    
	rows, err := pool.Query(ctx, query, request)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Requests, 0)

	for rows.Next() {
		var a Requests

		if err := rows.Scan(
			&a.First,
			&a.Last,
			&a.Id,
			&a.Previous,
			&a.Amount,
			&a.Total_Requested_Amount,
			&a.Details,
			&a.Status,
			&a.RegNum,
			&a.ParentStatus,
			&a.EmploymentStatus,
			&a.Programme,
		); err != nil {
			return nil, err
		}

		results = append(results, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetStatementRequests fetches unique student records from financial_history with status 'sent'
func GetStatementRequests(pool *pgxpool.Pool, ctx context.Context) ([]StatementRequest, error) {
	query := `
		SELECT DISTINCT ON (fh.student_id)
			u.name,
			u.surname,
			COALESCE(app.registration_number, ''),
			COALESCE(app.parent_guardian_status, ''),
			COALESCE(app.guardian_employment_status, ''),
			COALESCE(app.other_financial_support, 'N/A'),
			fh.request
		FROM financial_history fh
		JOIN users u ON u.user_id = fh.student_id
		LEFT JOIN applications app ON app.user_id = fh.student_id
		WHERE fh.request = 'sent'
		ORDER BY fh.student_id, fh.updated_at DESC;
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]StatementRequest, 0)
	for rows.Next() {
		var r StatementRequest
		if err := rows.Scan(
			&r.First,
			&r.Last,
			&r.RegNum,
			&r.ParentStatus,
			&r.EmploymentStatus,
			&r.Income,
			&r.Status,
		); err != nil {
			return nil, err
		}
		// Dummy priority for now
		r.Priority = "Low"
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func AcceptRequest(tx pgx.Tx, ctx context.Context, reqId int64) (string, error) {
	query := `UPDATE financial_request SET request_status = 'approved' WHERE request_id = $1`
	_, err := tx.Exec(ctx, query, reqId)
	if err != nil {
		return "", err
	}
	return "Request Submitted successfully", nil
}

func RejectRequest(tx pgx.Tx, ctx context.Context, reqId int64) (string, error) {
	query := `UPDATE financial_request SET request_status = 'rejected' WHERE request_id = $1`
	_, err := tx.Exec(ctx, query, reqId)
	if err != nil {
		return "", err
	}
	return "Request Rejected successfully", nil
}

func GetTotalAmount(
	pool *pgxpool.Pool,
	ctx context.Context,
) ([]float64, error) {

	query := `
	SELECT payment_amount
	FROM financial_request
	WHERE request_status = 'approved'::request;
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []float64

	for rows.Next() {
		var amount float64

		if err := rows.Scan(&amount); err != nil {
			return nil, err
		}

		results = append(results, amount)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}