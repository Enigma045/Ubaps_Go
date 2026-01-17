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

	semester, err := getFormValue(formData, "semester")
	if err != nil {
		return err
	}
	date, err := getFormValue(formData, "date")
	if err != nil {
		return err
	}
	detail, err := getFormValue(formData, "detail")
	if err != nil {
		return err
	}
	amount, err := getFormValue(formData, "amount")
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
