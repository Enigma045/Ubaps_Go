package Handles

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"ubaps/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Applicant struct {
	First  string `json:"first_name"`
	Last   string `json:"last_name"`
	Status string `json:"status"`
	Application_date utils.AutoTime `json:"application_date"`
	Dob utils.AutoTime `json:"dob"`
	Gender string `json:"gender"`
	Home_district string `json:"home_district"`
	Programme string `json:"programme"`
	Registration_number string `json:"registration_number"`
	Accommodation string `json:"accommodation"`
	Parent_guardian_status string `json:"parent_guardian_status"`
	Guardian_employment_status string `json:"guardian_employment_status"`
	Relative_support string `json:"relative_support"`
	Bursary_amount string `json:"bursary_amount"`
	Reason sql.NullString `json:"reason"`
	SchemeName string `json:"scheme_name"`
	Comments json.RawMessage `json:"comments"`
}

type ReportFilters struct {
	Statuses    []string `json:"statuses"`
	Search      string   `json:"search"`
	Department  []string `json:"department"`
	Scheme      []string `json:"scheme"`
	Parent      []string `json:"parent"`
	Employment  []string `json:"employment"`
	Gender      []string `json:"gender"`
	DateStart   string   `json:"date_start"`
	DateEnd     string   `json:"date_end"`
	CohortStart string   `json:"cohort_start"`
	CohortEnd   string   `json:"cohort_end"`
}

func Applicants(
	pool *pgxpool.Pool,
	ctx context.Context,
	filters ReportFilters,
	limit, offset int,
	year, semester, month string,
) ([]Applicant, int, error) {

	// Safety: avoid SQL error on empty slice
	if len(filters.Statuses) == 0 {
		return []Applicant{}, 0, nil
	}

	var realStatuses []string
	var feesPaidRequested bool
	for _, s := range filters.Statuses {
		realStatuses = append(realStatuses, s)
		if s == "paid" {
			feesPaidRequested = true
		}
	}

	if len(realStatuses) == 0 && !feesPaidRequested {
		return []Applicant{}, 0, nil
	}

	whereClause := "WHERE ((a.status = ANY($1)) OR ($2 = true AND EXISTS (SELECT 1 FROM financial_request fr WHERE fr.user_id = a.user_id AND fr.request_status = 'approved')))"
	params := []interface{}{realStatuses, feesPaidRequested}
	paramCount := 3

	if filters.Search != "" {
		whereClause += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.surname ILIKE $%d OR a.registration_number ILIKE $%d)", paramCount, paramCount, paramCount)
		params = append(params, "%"+filters.Search+"%")
		paramCount++
	}
	if len(filters.Department) > 0 {
		patterns := make([]string, len(filters.Department))
		for i, d := range filters.Department {
			patterns[i] = d + "%"
		}
		whereClause += fmt.Sprintf(" AND a.registration_number ILIKE ANY($%d)", paramCount)
		params = append(params, patterns)
		paramCount++
	}
	if len(filters.Scheme) > 0 {
		whereClause += fmt.Sprintf(" AND s.scheme_name = ANY($%d)", paramCount)
		params = append(params, filters.Scheme)
		paramCount++
	}
	if len(filters.Parent) > 0 {
		whereClause += fmt.Sprintf(" AND a.parent_guardian_status = ANY($%d)", paramCount)
		params = append(params, filters.Parent)
		paramCount++
	}
	if len(filters.Employment) > 0 {
		whereClause += fmt.Sprintf(" AND a.guardian_employment_status = ANY($%d)", paramCount)
		params = append(params, filters.Employment)
		paramCount++
	}
	if len(filters.Gender) > 0 {
		whereClause += fmt.Sprintf(" AND a.gender = ANY($%d)", paramCount)
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

	if filters.DateStart != "" {
		whereClause += fmt.Sprintf(" AND a.application_date >= $%d", paramCount)
		params = append(params, filters.DateStart)
		paramCount++
	}
	if filters.DateEnd != "" {
		whereClause += fmt.Sprintf(" AND a.application_date <= $%d", paramCount)
		params = append(params, filters.DateEnd+" 23:59:59")
		paramCount++
	}

	if filters.CohortStart != "" && filters.CohortEnd != "" {
		whereClause += fmt.Sprintf(" AND SUBSTRING(a.registration_number FROM '\\d+$')::integer BETWEEN $%d AND $%d", paramCount, paramCount+1)
		params = append(params, filters.CohortStart, filters.CohortEnd)
		paramCount += 2
	} else if filters.CohortStart != "" {
		whereClause += fmt.Sprintf(" AND SUBSTRING(a.registration_number FROM '\\d+$')::integer >= $%d", paramCount)
		params = append(params, filters.CohortStart)
		paramCount++
	} else if filters.CohortEnd != "" {
		whereClause += fmt.Sprintf(" AND SUBSTRING(a.registration_number FROM '\\d+$')::integer <= $%d", paramCount)
		params = append(params, filters.CohortEnd)
		paramCount++
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
    u.name,
    u.surname ,
	a.application_date,
	a.date_of_birth,
	a.status,
	a.gender,
	a.home_district,
	a.programme,
	a.registration_number,
	a.accommodation,
	a.parent_guardian_status,
	a.guardian_employment_status,
	a.other_financial_support,
	a.bursary_amount,
	a.reason_for_bursary,
	COALESCE(s.scheme_name, 'No Scheme')
    FROM applications a
    JOIN users u ON u.user_id = a.user_id
	LEFT JOIN bursary_schemes s ON s.scheme_id = a.scheme_id
    %s
    ORDER BY a.application_date DESC
    LIMIT $%d OFFSET $%d;
	`, whereClause, limitIdx, offsetIdx)

	rows, err := pool.Query(ctx, query, mainParams...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results := make([]Applicant, 0)
	for rows.Next() {
		var a Applicant
		if err := rows.Scan(
			&a.First,
			&a.Last,
			&a.Application_date,
			&a.Dob,
			&a.Status,
			&a.Gender,
			&a.Home_district,
			&a.Programme,
			&a.Registration_number,
			&a.Accommodation,
			&a.Parent_guardian_status,
			&a.Guardian_employment_status,
			&a.Relative_support,
			&a.Bursary_amount,
			&a.Reason,
			&a.SchemeName,
		); err != nil {
			return nil, 0, err
		}
		results = append(results, a)
	}
	return results, total, nil
}


func ConsiderStudent(tx pgx.Tx, ctx context.Context, userId int64,Role string) (string, error) {
//check if application is selected
	Check := `
    SELECT EXISTS (
        SELECT 1
        FROM applications
        WHERE user_id = $1
        AND status = 'selected'
    );
`

var exists bool

err := tx.QueryRow(ctx, Check, userId).Scan(&exists)
if err != nil {
    return "", err
}

if exists {
    return "Application already selected", nil
}
	
//check


	role := fmt.Sprintf("%s_approval_status",Role)
	query := fmt.Sprintf(`UPDATE applications SET %s = 'approved' , status = 'considering' WHERE user_id = $1`,role)
	_, err = tx.Exec(ctx, query, userId)
	if err != nil {
		return "", err
	}
	return "Application Considered successfully", nil
}

func RejectStudent(tx pgx.Tx, ctx context.Context, userId int64,Role string) (string, error) {
	//check if application is selected
	Check := `
    SELECT EXISTS (
        SELECT 1
        FROM applications
        WHERE user_id = $1
        AND status = 'selected'
    );
    `

   var exists bool

   err := tx.QueryRow(ctx, Check, userId).Scan(&exists)
   if err != nil {
    return "", err
   }

   if exists {
    return "Application already selected", nil
   }
	
   //check

   //check if application is rejected
	Check2 := `
    SELECT EXISTS (
        SELECT 1
        FROM applications
        WHERE user_id = $1
        AND status = 'not selected'
    );
    `

   var exists2 bool

   err = tx.QueryRow(ctx, Check2, userId).Scan(&exists2)
   if err != nil {
    return "", err
   }

   if exists2 {
    return "Application already rejected", nil
   }
	
   //check

	status := "considering"
	if Role == "registrar" {
		status = "not selected"
	}

	role := fmt.Sprintf("%s_approval_status", Role)
	query := fmt.Sprintf(`UPDATE applications SET %s = 'rejected' , status = $2 WHERE user_id = $1`, role)
	_, err = tx.Exec(ctx, query, userId, status)
	if err != nil {
		return "", err
	}
	return "Application Rejected successfully", nil
}

func PayInstallment(tx pgx.Tx, ctx context.Context, userId int64) (string, error) {
	// 1. Update application status
	fmt.Println("Processing payment for user:", userId)
	queryApp := `UPDATE applications SET status = 'paid' WHERE user_id = $1`
	_, err := tx.Exec(ctx, queryApp, userId)
	if err != nil {
		return "", err
	}

	// 2. Approve any pending financial requests for this user
	// queryReq := `UPDATE financial_request SET request_status = 'approved' WHERE user_id = $1 AND request_status = 'pending'`
	// _, err = tx.Exec(ctx, queryReq, userId) // We don't strictly care if no pending requests existed
    // if err != nil {
	// 	return "", err
	// }
	 
	return "Payment processed successfully", nil
}
func RollbackSelection(tx pgx.Tx, ctx context.Context, userId int64) (string, error) {
	// 1. Get scheme_id and bursary_amount
	var schemeId int64
	var amountStr string
	err := tx.QueryRow(ctx, "SELECT scheme_id, bursary_amount FROM applications WHERE user_id = $1", userId).Scan(&schemeId, &amountStr)
	if err != nil {
		return "", fmt.Errorf("failed to fetch application info: %w", err)
	}

	// 2. Return money to scheme
	amount, _ := utils.Strtofloat(amountStr)
	_, err = tx.Exec(ctx, "UPDATE bursary_schemes SET available_balance = available_balance + $1 WHERE scheme_id = $2", amount, schemeId)
	if err != nil {
		return "", fmt.Errorf("failed to restore scheme balance: %w", err)
	}

	// 3. Reset application
	_, err = tx.Exec(ctx, "UPDATE applications SET status = 'submitted', scheme_id = NULL, bursary_amount = '0' WHERE user_id = $1", userId)
	if err != nil {
		return "", fmt.Errorf("failed to reset application status: %w", err)
	}

	return "Selection rolled back successfully", nil
}

func AddComment(tx pgx.Tx, ctx context.Context, userId int64, comment string, userName string, userRole string) error {
	newComment := map[string]string{
		"name": userName,
		"role": userRole,
		"text": comment,
		"date": utils.Floattostr(float64(utils.AutoTime(context.Background()).Time().Unix())), // Using existing helpers or simple format
	}
	// Re-formatted date for readability
	newComment["date"] = utils.Floattostr(float64(time.Now().Unix())) 

	commentJSON, _ := json.Marshal(newComment)

	query := `
        UPDATE applications 
        SET comments = COALESCE(comments, '[]'::jsonb) || jsonb_build_array($1::jsonb) 
        WHERE user_id = $2
    `
	_, err := tx.Exec(ctx, query, string(commentJSON), userId)
	if err != nil {
		return fmt.Errorf("failed to add comment to DB: %w", err)
	}
	return nil
}
