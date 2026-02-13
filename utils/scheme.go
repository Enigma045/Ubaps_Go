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

type Scheme struct{
     Name string `json:"scheme_name"`
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
		totalFund,
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

func GetScheme(pool *pgxpool.Pool,ctx context.Context) ([]Scheme,error) {

	query := `SELECT scheme_name FROM bursary_schemes`

    rows,err := pool.Query(ctx,query)
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	var schemes []Scheme

	for rows.Next() {
		var scheme Scheme

		err := rows.Scan(
			&scheme.Name,
		)
		if err != nil {
			return nil,err
		}

		schemes = append(schemes,scheme)
	}

	if err := rows.Err(); err != nil {
		return nil,err
	}

	return schemes,nil
}

func GetSchemeId(schmename string,tx pgx.Tx,ctx context.Context)(int64,error){
    var SchemeID int64

	query := `SELECT scheme_id FROM bursary_schemes WHERE scheme_name = $1`

	err :=tx.QueryRow(ctx,query,schmename).Scan(&SchemeID)
    if err != nil {
		return 0,err
	}

	return SchemeID,nil
}

func CheckSchemeAmount(schmename string,tx pgx.Tx,ctx context.Context,amount string)(error){
    //var SchemeID int64

	query := `
	SELECT EXISTS(
	    SELECT 1 FROM bursary_schemes WHERE scheme_name = $1 AND available_balance > $2
	)
	`

	
//   check := `
//SELECT EXISTS(
//    SELECT 1 FROM bursary_schemes WHERE scheme_name = $1
//)
//`
   
    var exists bool

	err :=tx.QueryRow(ctx,query,schmename,amount).Scan(&exists)
    if err != nil {
		return err
	}

	if !exists {
        return fmt.Errorf("amount excceds scheme amount or scheme does not exist")
    }


	return nil
}

func SendScheme_Info(tx pgx.Tx,ctx context.Context,user_id int64,scheme_id int64,amount string) (error) {

	query := `UPDATE applications SET scheme_id=$1, status=$2, bursary_amount=$3 WHERE user_id=$4`

	_,err := tx.Exec(ctx,query,scheme_id,"selected",amount,user_id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateScheme_Amount(tx pgx.Tx,ctx context.Context,scheme_id int64,assighned_amount string,available_amount float64) (error) {

    value,err := Strtofloat(assighned_amount)
    if err != nil{
		log.Println(err)
		return err
	}

	deducted_amount := (available_amount-value)

	convamount := Floattostr(deducted_amount)
	 
	query := `UPDATE bursary_schemes SET available_balance=$1 WHERE  scheme_id=$2`

	_,err = tx.Exec(ctx,query,convamount,scheme_id,)
	if err != nil {
		return err
	}

	return nil
}

func GetAvailableAmount(tx pgx.Tx,ctx context.Context,scheme_id int64)(float64,error){
	var available_amount string

	query := `SELECT available_balance FROM bursary_schemes WHERE scheme_id=$1`

	err := tx.QueryRow(ctx,query,scheme_id).Scan(&available_amount)
	if err != nil {
		return 0,err
	}

	value,err := Strtofloat(available_amount)
	if err != nil {
		log.Println("Failed to convert to float")
		return 0,err
	}

	return value,nil
}