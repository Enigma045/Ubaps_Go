package Handles

import (
	"context"
	"database/sql"
	"fmt"
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

}

func Applicants(
	pool *pgxpool.Pool,
	ctx context.Context,
	applicants []string,
	limit, offset int,
) ([]Applicant, int, error) {

	// Safety: avoid SQL error on empty slice
	if len(applicants) == 0 {
		return []Applicant{}, 0, nil
	}

	// Filter out virtual statuses like 'paid' (if not in enum) to avoid Postgres Enum errors
	// Note: 'paid' is now in the enum, but we still keep the feesPaidRequested logic for financial_request check
	var realStatuses []string
	var feesPaidRequested bool
	for _, s := range applicants {
		realStatuses = append(realStatuses, s)
		if s == "paid" {
			feesPaidRequested = true
		}
	}

	// If no real statuses were provided and paid wasn't requested, return early
	if len(realStatuses) == 0 && !feesPaidRequested {
		return []Applicant{}, 0, nil
	}

	countQuery := `
		SELECT COUNT(*) 
		FROM applications a
		WHERE (a.status = ANY($1))
		   OR ($2 = true AND EXISTS (
		       SELECT 1 FROM financial_request fr 
		       WHERE fr.user_id = a.user_id AND fr.request_status = 'approved'
		   ))
	`
	var total int
	err := pool.QueryRow(ctx, countQuery, realStatuses, feesPaidRequested).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
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
    WHERE (a.status = ANY($1))
       OR ($4 = true AND EXISTS (
           SELECT 1 FROM financial_request fr 
           WHERE fr.user_id = a.user_id AND fr.request_status = 'approved'
       ))
    ORDER BY a.application_date DESC
    LIMIT $2 OFFSET $3;
	`

	rows, err := pool.Query(ctx, query, realStatuses, limit, offset, feesPaidRequested)
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
