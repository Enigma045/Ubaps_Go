package Handles

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FinanceReportResult struct {
	Stats map[string]interface{} `json:"stats"`
	Rows  [][]interface{}        `json:"rows"`
}

func GetDisbursementStatusReport(pool *pgxpool.Pool, ctx context.Context) (FinanceReportResult, error) {
	var result FinanceReportResult
	result.Stats = make(map[string]interface{})
	result.Rows = [][]interface{}{}

	// Stats
	queryStats := `
		SELECT 
			COALESCE(SUM(assigned_amount), 0),
			COALESCE(SUM(assigned_amount) FILTER (WHERE payed = true), 0),
			COALESCE(SUM(assigned_amount) FILTER (WHERE payed = false), 0)
		FROM dusbersment
	`
	var total, paid, pending float64
	err := pool.QueryRow(ctx, queryStats).Scan(&total, &paid, &pending)
	if err != nil {
		return result, err
	}
	result.Stats["total_assigned"] = total
	result.Stats["total_paid"] = paid
	result.Stats["total_pending"] = pending

	// Rows
	queryRows := `
		SELECT 
			u.name || ' ' || u.surname,
			a.registration_number,
			COALESCE(d.assigned_amount, 0),
			CASE WHEN d.payed = true THEN 'Paid' ELSE 'Unpaid' END,
			SUBSTRING(a.registration_number FROM '^[A-Za-z]+'),
			a.scheme_id
		FROM applications a
		JOIN users u ON a.user_id = u.user_id
		LEFT JOIN dusbersment d ON d.user_id = u.user_id
		WHERE a.status = 'selected' OR a.status = 'paid'
		ORDER BY a.registration_number ASC
	`
	rows, err := pool.Query(ctx, queryRows)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var name, reg, status, programme string
		var amount float64
		var schemeID int
		if err := rows.Scan(&name, &reg, &amount, &status, &programme, &schemeID); err != nil {
			continue
		}
		result.Rows = append(result.Rows, []interface{}{name, reg, fmt.Sprintf("MWK %.2f", amount), status, programme, fmt.Sprintf("SCH-%02d", schemeID)})
	}

	return result, nil
}

func GetFundUtilisationReport(pool *pgxpool.Pool, ctx context.Context) (FinanceReportResult, error) {
	var result FinanceReportResult
	result.Stats = make(map[string]interface{})
	result.Rows = [][]interface{}{}

	// Stats
	queryStats := `
		SELECT 
			COUNT(*),
			COALESCE(SUM(total_fund_amount - available_balance), 0),
			COALESCE(SUM(available_balance), 0)
		FROM bursary_schemes
	`
	var count int
	var committed, available float64
	err := pool.QueryRow(ctx, queryStats).Scan(&count, &committed, &available)
	if err != nil {
		return result, err
	}
	result.Stats["scheme_count"] = count
	result.Stats["total_committed"] = committed
	result.Stats["net_available"] = available

	// Rows
	queryRows := `
		SELECT 
			scheme_name,
			total_fund_amount,
			available_balance,
			benefactor_name,
			CASE WHEN is_active = true THEN 'Active' ELSE 'Inactive' END
		FROM bursary_schemes
		ORDER BY scheme_name ASC
	`
	rows, err := pool.Query(ctx, queryRows)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var name, benefactor, active string
		var total, avail float64
		if err := rows.Scan(&name, &total, &avail, &benefactor, &active); err != nil {
			continue
		}
		result.Rows = append(result.Rows, []interface{}{name, fmt.Sprintf("MWK %.2f", total), fmt.Sprintf("MWK %.2f", avail), benefactor, active})
	}

	return result, nil
}

func GetPaymentHistoryReport(pool *pgxpool.Pool, ctx context.Context) (FinanceReportResult, error) {
	var result FinanceReportResult
	result.Stats = make(map[string]interface{})
	result.Rows = [][]interface{}{}

	// Stats
	queryStats := `
		SELECT 
			COUNT(*),
			COALESCE(SUM(payment_amount), 0)
		FROM financial_history
		WHERE request = 'answered'
	`
	var count int
	var volume float64
	err := pool.QueryRow(ctx, queryStats).Scan(&count, &volume)
	if err != nil {
		return result, err
	}
	result.Stats["tx_count"] = count
	result.Stats["total_volume"] = volume

	// Rows
	queryRows := `
		SELECT 
			COALESCE(TO_CHAR(payment_date, 'YYYY-MM-DD'), 'N/A'),
			payment_amount,
			semester,
			cumulative_amount,
			request,
			CASE WHEN installment_status = true THEN 'Full' ELSE 'Partial' END
		FROM financial_history
		ORDER BY payment_date DESC
	`
	rows, err := pool.Query(ctx, queryRows)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var date, sem, status, inst string
		var amount, cum float64
		if err := rows.Scan(&date, &amount, &sem, &cum, &status, &inst); err != nil {
			continue
		}
		result.Rows = append(result.Rows, []interface{}{date, fmt.Sprintf("MWK %.2f", amount), sem, fmt.Sprintf("MWK %.2f", cum), status, inst})
	}

	return result, nil
}

func GetPendingDisbursementsReport(pool *pgxpool.Pool, ctx context.Context) (FinanceReportResult, error) {
	var result FinanceReportResult
	result.Stats = make(map[string]interface{})
	result.Rows = [][]interface{}{}

	// Stats
	queryStats := `
		SELECT 
			COUNT(*),
			COALESCE(SUM(assigned_amount), 0)
		FROM dusbersment
		WHERE payed = false
	`
	var count int
	var capital float64
	err := pool.QueryRow(ctx, queryStats).Scan(&count, &capital)
	if err != nil {
		return result, err
	}
	result.Stats["pending_count"] = count
	result.Stats["required_capital"] = capital

	// Rows
	queryRows := `
		SELECT 
			a.registration_number,
			d.assigned_amount,
			SUBSTRING(a.registration_number FROM '^[A-Za-z]+'),
			COALESCE(TO_CHAR(a.application_date, 'YYYY-MM-DD HH24:MI'), 'N/A')
		FROM applications a
		JOIN dusbersment d ON d.user_id = a.user_id
		WHERE d.payed = false AND (a.status = 'selected' OR a.status = 'paid')
		ORDER BY a.application_date ASC
	`
	rows, err := pool.Query(ctx, queryRows)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var reg, programme, approvedAt string
		var amount float64
		if err := rows.Scan(&reg, &amount, &programme, &approvedAt); err != nil {
			continue
		}
		result.Rows = append(result.Rows, []interface{}{reg, fmt.Sprintf("MWK %.2f", amount), programme, approvedAt, "Process Now"})
	}

	return result, nil
}
