package Handles

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ApprovalStatus struct {
	RegistrarApproval     sql.NullString `json:"registrar_approval_status"`
	DeanStudentApproval   sql.NullString `json:"dean_of_student_approval_status"`
	DeanFacultApproval    sql.NullString `json:"dean_of_facult_approval_status"`
	DeanScienceApproval   sql.NullString `json:"dean_of_science_approval_status"`
	FinanceOfficeApproval sql.NullString `json:"finance_office_approval_status"`
	Status                string         `json:"status"`
	BursaryAmount         string         `json:"bursary_amount"`
	SchemeName            sql.NullString `json:"scheme_name"`
}

func GetApplicationStatus(
	pool *pgxpool.Pool,
	ctx context.Context,
	registrationNumber string,
) (ApprovalStatus, error) {

	var status ApprovalStatus

	query := `
		SELECT 
			a.registrar_approval_status,
			a.dean_of_student_approval_status,
			a.dean_of_facult_approval_status,
			a.dean_of_science_approval_status,
			a.finance_office_approval_status,
			a.status,
			a.bursary_amount,
			b.scheme_name
		FROM applications a
		LEFT JOIN bursary_schemes b ON a.scheme_id = b.scheme_id
		WHERE a.registration_number = $1
	`

	err := pool.QueryRow(ctx, query, registrationNumber).Scan(
		&status.RegistrarApproval,
		&status.DeanStudentApproval,
		&status.DeanFacultApproval,
		&status.DeanScienceApproval,
		&status.FinanceOfficeApproval,
		&status.Status,
		&status.BursaryAmount,
		&status.SchemeName,
	)
	if err != nil {
		return ApprovalStatus{}, err
	}

	return status, nil
}
