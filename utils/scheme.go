package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Benefactor struct {
    Scheme_name string `json:"scheme_name"`
	Benefactor_name string `json:"benefactor_name"`
	Benefactor_email string `json:"benefactor_email"`
	Total_fund_amount float64 `json:"total_fund_amount"`
	Available_balance float64 `json:"available_balance"`
	Gender_restriction string `json:"gender_restriction"`
	Conditions string `json:"conditions"`
}

func Scheme_Operations(
	tx pgx.Tx,
	ctx context.Context,
	formData map[string][]string,
	userid int64,
) error {

	schemeName, err := GetFormValue(formData, "scheme_name")
	if err != nil {
		return err
	}
	benefactorName, err := GetFormValue(formData, "benefactor_name")
	if err != nil {
		return err
	}
	benefactorEmail, err := GetFormValue(formData, "benefactor_email")
	if err != nil {
		return err
	}
	totalFund, err := GetFormValue(formData, "total_fund_amount")
	if err != nil {
		return err
	}
    
	// optional / default fields

	genderRestriction := "both"
	if v, err := GetFormValue(formData, "gender_restriction"); err == nil {
		genderRestriction = v
	}

	conditions := ""
	if v, err := GetFormValue(formData, "conditions"); err == nil {
		conditions = v
	}

	//check
    check := `
    SELECT EXISTS(
        SELECT 1 FROM bursary_schemes WHERE scheme_name = $1
    )
`

    var exists bool

    err = tx.QueryRow(ctx, check, schemeName).Scan(&exists)
    if err != nil {
        return fmt.Errorf("db check failed: %w", err)
    }

    if exists {
        return fmt.Errorf("scheme already exists")
    }

	//

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



func GetBenefactor(
	pool *pgxpool.Pool,
	 ctx context.Context,
) ([]Benefactor, error) {

	query := `
		SELECT 
		scheme_name,
		benefactor_name,
		benefactor_email,
		total_fund_amount,
		available_balance,
		gender_restriction,
		conditions
		FROM bursary_schemes
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ben []Benefactor

	for rows.Next() {
		var benify Benefactor

		err := rows.Scan(
			&benify.Scheme_name,
			&benify.Benefactor_name,
			&benify.Benefactor_email,
			&benify.Total_fund_amount,
			&benify.Available_balance,
			&benify.Gender_restriction,
			&benify.Conditions,
		)
		if err != nil {
			return nil, err
		}

		ben = append(ben, benify)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ben, nil
}

func DeleteBenefactor(tx pgx.Tx,
	            ctx context.Context,
				 emailreq string) (error) {

	query := `DELETE FROM
	bursary_schemes WHERE
	scheme_name = $1`

	_,err := tx.Exec(ctx,query,emailreq)
	if err != nil {
    return err
	}

	return nil
}