package Handles

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ComprehensiveReportRow represents a single line in our comprehensive report.
type ComprehensiveReportRow struct {
	ApplicationID int     `json:"application_id"`
	UserID        int     `json:"user_id"`
	Name          string  `json:"name"`
	Surname       string  `json:"surname"`
	Email         string  `json:"email"`
	Status        string  `json:"status"`
	SchemeName    string  `json:"scheme_name"`
	BursaryAmount float64 `json:"bursary_amount"`
	HomeDistrict  string  `json:"home_district"`
	Gender        string  `json:"gender"`
	PaymentAmount float64 `json:"payment_amount"`
	RequestStatus string  `json:"request_status"`
	AppliedAt     string  `json:"applied_at"`
}

// GetComprehensiveReport queries multiple tables to form a detailed report
func GetComprehensiveReport(pool *pgxpool.Pool, ctx context.Context, limit, offset int) ([]ComprehensiveReportRow, int, error) {
	countQuery := `SELECT COUNT(*) FROM applications`
	var total int
	err := pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			a.application_id, 
			u.user_id, 
			COALESCE(u.name, ''), 
			COALESCE(u.surname, ''), 
			u.email, 
			a.status, 
			COALESCE(s.scheme_name, 'No Scheme'),
			COALESCE(a.bursary_amount, 0),
			COALESCE(a.home_district, ''),
			COALESCE(a.gender, ''),
			COALESCE(fr.payment_amount, 0),
			COALESCE(fr.request_status, 'No Request'),
			COALESCE(TO_CHAR(a.created_at, 'YYYY-MM-DD HH24:MI:SS'), '')
		FROM applications a
		JOIN users u ON u.user_id = a.user_id
		LEFT JOIN bursary_schemes s ON s.scheme_id = a.scheme_id
		LEFT JOIN financial_request fr ON fr.user_id = u.user_id
		ORDER BY a.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println("Error executing report query:", err)
		return nil, 0, err
	}
	defer rows.Close()

	var reportData []ComprehensiveReportRow
	for rows.Next() {
		var row ComprehensiveReportRow
		err := rows.Scan(
			&row.ApplicationID,
			&row.UserID,
			&row.Name,
			&row.Surname,
			&row.Email,
			&row.Status,
			&row.SchemeName,
			&row.BursaryAmount,
			&row.HomeDistrict,
			&row.Gender,
			&row.PaymentAmount,
			&row.RequestStatus,
			&row.AppliedAt,
		)
		if err != nil {
			log.Println("Error scanning report row:", err)
			continue
		}
		reportData = append(reportData, row)
	}

	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	return reportData, total, nil
}
