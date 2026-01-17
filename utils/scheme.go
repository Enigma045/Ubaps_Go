package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func Scheme_Operations(
	tx pgx.Tx,
	ctx context.Context,
	formData map[string][]string,
	userid int64,
) error {

	schemeName, err := getFormValue(formData, "scheme_name")
	if err != nil {
		return err
	}
	benefactorName, err := getFormValue(formData, "benefactor_name")
	if err != nil {
		return err
	}
	benefactorEmail, err := getFormValue(formData, "benefactor_email")
	if err != nil {
		return err
	}
	totalFund, err := getFormValue(formData, "total_fund_amount")
	if err != nil {
		return err
	}

	// optional / default fields

	genderRestriction := ""
	if v, err := getFormValue(formData, "gender_restriction"); err == nil {
		genderRestriction = v
	}

	conditions := ""
	if v, err := getFormValue(formData, "conditions"); err == nil {
		conditions = v
	}

	query := `
		INSERT INTO bursary_schemes (
			scheme_name,
			benefactor_name,
			benefactor_email,
			total_fund_amount,
			available_balance,
			gender_restriction,
			conditions,
			user_id
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8);
	`

	row, err := tx.Exec(
		ctx,
		query,
		schemeName,
		benefactorName,
		benefactorEmail,
		totalFund,
		20000,
		genderRestriction,
		conditions,
		userid,
	)
	if err != nil {
		return fmt.Errorf("db insert failed: %w", err)
	}

	log.Println("Inserted rows:", row.RowsAffected())
	return nil
}
