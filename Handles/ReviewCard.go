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
}

func GetApplicationStatus(
	pool *pgxpool.Pool,
	ctx context.Context,
	registrationNumber string,
) (ApprovalStatus, error) {

	var status ApprovalStatus

	query := `
		SELECT 
			registrar_approval_status,
			dean_of_student_approval_status,
			dean_of_facult_approval_status,
			dean_of_science_approval_status,
			finance_office_approval_status
		FROM applications
		WHERE registration_number = $1
	`

	err := pool.QueryRow(ctx, query, registrationNumber).Scan(
		&status.RegistrarApproval,
		&status.DeanStudentApproval,
		&status.DeanFacultApproval,
		&status.DeanScienceApproval,
		&status.FinanceOfficeApproval,
	)
	if err != nil {
		return ApprovalStatus{}, err
	}

	return status, nil
}
