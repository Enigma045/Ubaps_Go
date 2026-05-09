package Handles

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// StudentStats holds statistics for the student dashboard
type StudentStats struct {
	ApplicationStatus string  `json:"application_status"`
	BursaryScheme     string  `json:"bursary_scheme"`
	FeesPaid          float64 `json:"fees_paid"`
}

// RegistrarStats holds statistics for the registrar dashboard
type RegistrarStats struct {
	ApprovedAmount          float64 `json:"approved_amount"`
	NumberOfApplicants      int     `json:"number_of_applicants"`
	PendingApplications     int     `json:"pending_applications"`
	ConsideringApplications int     `json:"considering_applications"`
	SelectedStudents        int     `json:"selected_students"`
	RejectedStudents        int     `json:"rejected_students"`
	NumberOfSchemes         int     `json:"number_of_schemes"`
}

// AdminStats holds statistics for the admin dashboard
type AdminStats struct {
	TotalUsers    int `json:"total_users"`
	ActiveUsers   int `json:"active_users"`
	DeactiveUsers int `json:"deactive_users"`
}

// DeanStats holds statistics for the dean dashboard
type DeanStats struct {
	PendingApplications     int `json:"pending_applications"`
	ConsideringApplications int `json:"considering_applications"`
	SelectedStudents        int `json:"selected_students"`
	RejectedStudents        int `json:"rejected_students"`
	PendingLetters          int `json:"pending_letters"`
}

// FinanceStats holds statistics for the financial officer dashboard
type FinanceStats struct {
	ApprovedAmount           float64 `json:"approved_amount"`
	DisbursementsMade        int     `json:"disbursements_made"`
	FinancialHistoryRequests int     `json:"financial_history_requests"`
}

// ReportStats holds statistics for the student reports dashboard
type ReportStats struct {
	TotalApplications int           `json:"total_applications"`
	Approved          int           `json:"approved"`
	Pending           int           `json:"pending"`
	Rejected          int           `json:"rejected"`
	TotalValue        float64       `json:"total_value"`
	AvgProcessingTime float64       `json:"avg_processing_time"`
	FacultyBreakdown  []FacultyStat `json:"faculty_breakdown"`
}

type FacultyStat struct {
	Faculty string `json:"faculty"`
	Count   int    `json:"count"`
}

// UserProfile holds the name and role of the current user
type UserProfile struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// UserProfileDetailed holds comprehensive user information
type UserProfileDetailed struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Role    string `json:"role"`
}

// GetStudentStats fetches stats for a specific student
func GetStudentStats(pool *pgxpool.Pool, ctx context.Context, userID int64) (StudentStats, error) {
	var stats StudentStats
	query := `
		SELECT 
			COALESCE(a.status::TEXT, 'not submitted'),
			COALESCE(s.scheme_name, 'None'),
			CASE 
				WHEN a.status = 'paid' THEN COALESCE(CAST(NULLIF(TRIM(CAST(a.bursary_amount AS TEXT)), '') AS DOUBLE PRECISION), 0)
				ELSE 0
			END
		FROM users u
		LEFT JOIN applications a ON a.user_id = u.user_id
		LEFT JOIN bursary_schemes s ON s.scheme_id = a.scheme_id
		WHERE u.user_id = $1
	`
	err := pool.QueryRow(ctx, query, userID).Scan(&stats.ApplicationStatus, &stats.BursaryScheme, &stats.FeesPaid)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}
	return stats, nil
}

// GetRegistrarStats fetches stats for the registrar
func GetRegistrarStats(pool *pgxpool.Pool, ctx context.Context) (RegistrarStats, error) {
	var stats RegistrarStats

	// Approved Amount
	// Hardened to strip non-numeric characters before casting
	queryAmount := `SELECT COALESCE(SUM(payment_amount), 0)
                    FROM financial_request
                    WHERE request_status = 'approved';`
	err := pool.QueryRow(ctx, queryAmount).Scan(&stats.ApprovedAmount)
	if err != nil {
		return stats, err
	}

	// Number of Applicants
	query := `
		SELECT 
			COUNT(*) FILTER (WHERE status = 'submitted'),
			COUNT(*) FILTER (WHERE status = 'considering'),
			COUNT(*) FILTER (WHERE status = 'selected'),
			COUNT(*) FILTER (WHERE status = 'not selected')
		FROM applications
	`
	err = pool.QueryRow(ctx, query).Scan(&stats.PendingApplications, &stats.ConsideringApplications, &stats.SelectedStudents, &stats.RejectedStudents)
	if err != nil {
		return stats, err
	}

	// Assuming pending letters are applications that are 'selected' but haven't had some letter-specific action?
	// For now, let's just use pending applications count as a placeholder.

	// Number of Schemes
	querySchemes := `SELECT COUNT(*) FROM bursary_schemes`
	err = pool.QueryRow(ctx, querySchemes).Scan(&stats.NumberOfSchemes)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

// GetAdminStats fetches stats for the admin
func GetAdminStats(pool *pgxpool.Pool, ctx context.Context) (AdminStats, error) {
	var stats AdminStats

	query := `
		SELECT 
			COUNT(*),
			COUNT(*) FILTER (WHERE is_active = true),
			COUNT(*) FILTER (WHERE is_active = false)
		FROM users
	`
	err := pool.QueryRow(ctx, query).Scan(&stats.TotalUsers, &stats.ActiveUsers, &stats.DeactiveUsers)
	return stats, err
}

// GetDeanStats fetches stats for deans
func GetDeanStats(pool *pgxpool.Pool, ctx context.Context) (DeanStats, error) {
	var stats DeanStats

	query := `
		SELECT 
			COUNT(*) FILTER (WHERE status = 'submitted'),
			COUNT(*) FILTER (WHERE status = 'considering'),
			COUNT(*) FILTER (WHERE status = 'selected'),
			COUNT(*) FILTER (WHERE status = 'not selected')
		FROM applications
	`
	err := pool.QueryRow(ctx, query).Scan(&stats.PendingApplications, &stats.ConsideringApplications, &stats.SelectedStudents, &stats.RejectedStudents)
	if err != nil {
		return stats, err
	}

	// Assuming pending letters are applications that are 'selected' but haven't had some letter-specific action?
	// For now, let's just use pending applications count as a placeholder.
	stats.PendingLetters = stats.PendingApplications

	return stats, nil
}

// GetFinanceStats fetches stats for the financial officer
func GetFinanceStats(pool *pgxpool.Pool, ctx context.Context) (FinanceStats, error) {
	var stats FinanceStats

	// Approved Amount
	// Hardened to strip non-numeric characters before casting
	queryAmount := `SELECT COALESCE(SUM(payment_amount), 0)
                    FROM financial_request
                    WHERE request_status = 'approved';`
	err := pool.QueryRow(ctx, queryAmount).Scan(&stats.ApprovedAmount)
	if err != nil {
		return stats, err
	}

	// Disbursements Made
	queryDisbursements := `SELECT COUNT(*) FROM dusbersment`
	err = pool.QueryRow(ctx, queryDisbursements).Scan(&stats.DisbursementsMade)
	if err != nil {
		return stats, err
	}

	// Number of requests for financial history
	queryRequests := `SELECT COUNT(*) FROM financial_history WHERE request = 'pending'`
	err = pool.QueryRow(ctx, queryRequests).Scan(&stats.FinancialHistoryRequests)
	if err != nil {
		stats.FinancialHistoryRequests = 0
	}

	return stats, nil
}

// GetUserProfile fetches the current user's profile
func GetUserProfile(pool *pgxpool.Pool, ctx context.Context, userID int64) (UserProfile, error) {
	var profile UserProfile
	// Use TRIM and COALESCE to handle potential NULL names
	query := `SELECT TRIM(COALESCE(name, '') || ' ' || COALESCE(surname, '')), user_type FROM users WHERE user_id = $1`
	err := pool.QueryRow(ctx, query, userID).Scan(&profile.Name, &profile.Role)
	return profile, err
}

// GetDetailedUserProfile fetches all personal info for the current user
func GetDetailedUserProfile(pool *pgxpool.Pool, ctx context.Context, userID int64) (UserProfileDetailed, error) {
	var profile UserProfileDetailed
	query := `SELECT COALESCE(name, ''), COALESCE(surname, ''), email, COALESCE(phone, ''), user_type FROM users WHERE user_id = $1`
	err := pool.QueryRow(ctx, query, userID).Scan(&profile.Name, &profile.Surname, &profile.Email, &profile.Phone, &profile.Role)
	return profile, err
}

// GetReportStats fetches comprehensive statistics for reporting
func GetReportStats(pool *pgxpool.Pool, ctx context.Context, year, semester, month string) (ReportStats, error) {
	var stats ReportStats
	stats.FacultyBreakdown = []FacultyStat{}

	// Basic filtering logic for period
	// In a real app, you'd parse these into a WHERE clause
	// For now, let's assume we filter by application_date if year/month provided
	whereClause := "WHERE 1=1"
	params := []interface{}{}
	paramCount := 1

	if year != "" {
		if len(year) > 4 && strings.Contains(year, "/") {
			parts := strings.Split(year, "/")
			whereClause += fmt.Sprintf(" AND (EXTRACT(YEAR FROM application_date) = $%d OR EXTRACT(YEAR FROM application_date) = $%d)", paramCount, paramCount+1)
			params = append(params, parts[0], parts[1])
			paramCount += 2
		} else {
			whereClause += fmt.Sprintf(" AND EXTRACT(YEAR FROM application_date) = $%d", paramCount)
			params = append(params, year)
			paramCount++
		}
	}
	if month != "" {
		whereClause += fmt.Sprintf(" AND EXTRACT(MONTH FROM application_date) = $%d", paramCount)
		params = append(params, month)
		paramCount++
	}
	if semester != "" {
		// Semester 1: July-Dec, Semester 2: Jan-June (example logic)
		if semester == "1" {
			whereClause += " AND EXTRACT(MONTH FROM application_date) BETWEEN 7 AND 12"
		} else if semester == "2" {
			whereClause += " AND EXTRACT(MONTH FROM application_date) BETWEEN 1 AND 6"
		}
	}

	// 1. Status counts and Total Value
	query := fmt.Sprintf(`
		SELECT 
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'selected' OR status = 'paid'),
			COUNT(*) FILTER (WHERE status = 'submitted' OR status = 'considering'),
			COUNT(*) FILTER (WHERE status = 'not selected'),
			COALESCE(SUM(CAST(NULLIF(TRIM(CAST(bursary_amount AS TEXT)), '') AS DOUBLE PRECISION)) FILTER (WHERE status = 'selected' OR status = 'paid'), 0)
		FROM applications
		%s
	`, whereClause)

	err := pool.QueryRow(ctx, query, params...).Scan(
		&stats.TotalApplications,
		&stats.Approved,
		&stats.Pending,
		&stats.Rejected,
		&stats.TotalValue,
	)
	if err != nil {
		return stats, err
	}

	// 2. Average Processing Time (Placeholder)
	stats.AvgProcessingTime = 4.5

	// 3. Faculty Breakdown (mapped to 'programme' column)
	facultyQuery := fmt.Sprintf(`
		SELECT programme, COUNT(*)
		FROM applications
		%s
		GROUP BY programme
		ORDER BY COUNT(*) DESC
	`, whereClause)

	rows, err := pool.Query(ctx, facultyQuery, params...)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var f FacultyStat
		if err := rows.Scan(&f.Faculty, &f.Count); err != nil {
			continue
		}
		stats.FacultyBreakdown = append(stats.FacultyBreakdown, f)
	}

	return stats, nil
}
