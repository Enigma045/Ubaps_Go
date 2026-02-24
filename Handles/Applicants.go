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

}

func Applicants(
	pool *pgxpool.Pool,
	ctx context.Context,
	applicants []string,
) ([]Applicant, error) {

	// Safety: avoid SQL error on empty slice
	if len(applicants) == 0 {
		return []Applicant{}, nil
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
	a.reason_for_bursary
    FROM applications a
    JOIN users u ON u.user_id = a.user_id
    WHERE a.status = ANY($1);

	`

	rows, err := pool.Query(ctx, query, applicants)
	if err != nil {
		return nil, err
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

	role := fmt.Sprintf("%s_approval_status",Role)
	query := fmt.Sprintf(`UPDATE applications SET %s = 'rejected' , status = 'considering' WHERE user_id = $1`,role)
	_, err = tx.Exec(ctx, query, userId)
	if err != nil {
		return "", err
	}
	return "Application Rejected successfully", nil
}
