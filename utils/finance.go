package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func Finance_Operations(
	tx pgx.Tx,
	ctx context.Context,
	formData map[string][]string,
	student_id int64,
) error {

	semester, err := GetFormValue(formData, "semester")
	if err != nil {
		return err
	}
	date, err := GetFormValue(formData, "date")
	if err != nil {
		return err
	}
	detail, err := GetFormValue(formData, "detail")
	if err != nil {
		return err
	}
	amount, err := GetFormValue(formData, "amount")
	if err != nil {
		return err
	}

	// optional / default fields

	// student_id := ""
	// if v, err := getFormValue(formData, "student_id"); err == nil {
	// 	student_id = v
	// }

	query := `
		INSERT INTO financial_history (
			semester,
			payment_date,
			details,
			payment_amount,
			student_id,
			request,
			updated_at,
			full_installment
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8);
	`

	row, err := tx.Exec(
		ctx,
		query,
		semester,
		date,
		detail,
		amount,
		student_id,
		"answered",
		time.Now(),
		891000,
	)
	if err != nil {
		return fmt.Errorf("db insert failed: %w", err)
	}

	log.Println("Inserted rows:", row.RowsAffected())
	return nil
}

// Request_Finance_Statement inserts a new record with status 'sent' into financial_history
func Request_Finance_Statement(tx pgx.Tx, ctx context.Context, studentID int64) error {
	query := `
		INSERT INTO financial_history (
			semester,
			payment_date,
			details,
			payment_amount,
			student_id,
			request,
			updated_at,
			full_installment
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8);
	`

	_, err := tx.Exec(
		ctx,
		query,
		"N/A",                // semester
		time.Now(),           // payment_date
		"Statement Requested", // details
		0,                    // payment_amount
		studentID,            // student_id
		"sent",               // request status
		time.Now(),           // updated_at
		0,                    // full_installment
	)
	if err != nil {
		return fmt.Errorf("failed to insert statement request: %w", err)
	}

	return nil
}
