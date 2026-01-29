package Handles

import (
	"context"
	"log"
	"ubaps/utils"

	"github.com/jackc/pgx/v5"
)

func Db_operation(tx pgx.Tx,ctx context.Context,userid int64,pamount int64,formData map[string][]string,) error {
	

	ramount, err := utils.GetFormValue(formData, "amount")
	if err != nil {
		return err
	}

	details, err := utils.GetFormValue(formData, "detail")
	if err != nil {
		return err
	}

	row,err := tx.Exec(ctx,
		`INSERT INTO financial_request (user_id,
		 previous_amount,
		 payment_amount,
		 details,
		 request_answered,
		 created_at,
		 updated_at
	)
		 VALUES ($1,$2,$3,$4,false,NOW(),NOW())`,
		 userid,pamount,ramount,details);

	if err != nil {
		log.Println("Error inserting financial request:", err)
     		return err
	}
    log.Println("Inserted rows:", row.RowsAffected())
	return nil
	
}