package Handles

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ComprehensiveReportRow represents a single line in our comprehensive report.
type ComprehensiveReportRow struct {
	ApplicationID      int     `json:"application_id"`
	UserID             int     `json:"user_id"`
	Name               string  `json:"name"`
	Surname            string  `json:"surname"`
	Email              string  `json:"email"`
	Status             string  `json:"status"`
	RegistrationNumber string  `json:"registration_number"`
	Programme          string  `json:"programme"`
	SchemeName         string  `json:"scheme_name"`
	BursaryAmount      float64 `json:"bursary_amount"`
	HomeDistrict       string  `json:"home_district"`
	Gender             string  `json:"gender"`
	PaymentAmount      float64 `json:"payment_amount"`
	RequestStatus      string  `json:"request_status"`
	AppliedAt          string  `json:"applied_at"`
}

// GetComprehensiveReport queries multiple tables to form a detailed report
func GetComprehensiveReport(pool *pgxpool.Pool, ctx context.Context, limit, offset int, year, semester, month string, filters ReportFilters) ([]ComprehensiveReportRow, int, error) {
	whereClause := "WHERE 1=1"
	params := []interface{}{}
	paramCount := 1

	if len(filters.Statuses) > 0 {
		var realStatuses []string
		var feesPaidRequested bool
		for _, s := range filters.Statuses {
			realStatuses = append(realStatuses, s)
			if s == "paid" {
				feesPaidRequested = true
			}
		}

		if len(realStatuses) > 0 || feesPaidRequested {
			whereClause += fmt.Sprintf(" AND ((a.status = ANY($%d)) OR ($%d = true AND EXISTS (SELECT 1 FROM financial_request fr WHERE fr.user_id = a.user_id AND fr.request_status = 'approved')))", paramCount, paramCount+1)
			params = append(params, realStatuses, feesPaidRequested)
			paramCount += 2
		}
	}

	if filters.Search != "" {
		whereClause += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.surname ILIKE $%d OR a.registration_number ILIKE $%d)", paramCount, paramCount, paramCount)
		params = append(params, "%"+filters.Search+"%")
		paramCount++
	}
	if filters.Department != "" {
		whereClause += fmt.Sprintf(" AND a.registration_number ILIKE $%d", paramCount)
		params = append(params, filters.Department+"%")
		paramCount++
	}
	if filters.Scheme != "" {
		whereClause += fmt.Sprintf(" AND s.scheme_name = $%d", paramCount)
		params = append(params, filters.Scheme)
		paramCount++
	}
	if filters.Parent != "" {
		whereClause += fmt.Sprintf(" AND a.parent_guardian_status = $%d", paramCount)
		params = append(params, filters.Parent)
		paramCount++
	}
	if filters.Employment != "" {
		whereClause += fmt.Sprintf(" AND a.guardian_employment_status = $%d", paramCount)
		params = append(params, filters.Employment)
		paramCount++
	}
	if filters.Gender != "" {
		whereClause += fmt.Sprintf(" AND a.gender = $%d", paramCount)
		params = append(params, filters.Gender)
		paramCount++
	}

	if year != "" {
		if len(year) > 4 && strings.Contains(year, "/") {
			parts := strings.Split(year, "/")
			whereClause += fmt.Sprintf(" AND (EXTRACT(YEAR FROM a.application_date) = $%d OR EXTRACT(YEAR FROM a.application_date) = $%d)", paramCount, paramCount+1)
			params = append(params, parts[0], parts[1])
			paramCount += 2
		} else {
			whereClause += fmt.Sprintf(" AND EXTRACT(YEAR FROM a.application_date) = $%d", paramCount)
			params = append(params, year)
			paramCount++
		}
	}
	if month != "" {
		whereClause += fmt.Sprintf(" AND EXTRACT(MONTH FROM a.application_date) = $%d", paramCount)
		params = append(params, month)
		paramCount++
	}
	if semester != "" {
		if semester == "1" {
			whereClause += " AND EXTRACT(MONTH FROM a.application_date) BETWEEN 7 AND 12"
		} else if semester == "2" {
			whereClause += " AND EXTRACT(MONTH FROM a.application_date) BETWEEN 1 AND 6"
		}
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM applications a
		JOIN users u ON u.user_id = a.user_id
		LEFT JOIN bursary_schemes s ON s.scheme_id = a.scheme_id
		%s
	`, whereClause)
	var total int
	err := pool.QueryRow(ctx, countQuery, params...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add limit and offset to params for the main query
	mainParams := append(params, limit, offset)
	limitIdx := len(mainParams) - 1
	offsetIdx := len(mainParams)

	query := fmt.Sprintf(`
		SELECT 
			a.application_id, 
			u.user_id, 
			COALESCE(u.name, ''), 
			COALESCE(u.surname, ''), 
			u.email, 
			a.status, 
			COALESCE(a.registration_number, ''),
			COALESCE(a.programme, ''),
			COALESCE(s.scheme_name, 'No Scheme'),
			COALESCE(a.bursary_amount, 0),
			COALESCE(a.home_district, ''),
			COALESCE(a.gender, ''),
			COALESCE(fr.payment_amount, 0),
			COALESCE(fr.request_status::TEXT, 'No Request'),
			COALESCE(TO_CHAR(a.created_at, 'YYYY-MM-DD HH24:MI:SS'), '')
		FROM applications a
		JOIN users u ON u.user_id = a.user_id
		LEFT JOIN bursary_schemes s ON s.scheme_id = a.scheme_id
		LEFT JOIN financial_request fr ON fr.user_id = u.user_id
		%s
		ORDER BY a.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, limitIdx, offsetIdx)

	rows, err := pool.Query(ctx, query, mainParams...)
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
			&row.RegistrationNumber,
			&row.Programme,
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
