package utils

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

// FirstFill ensures that a student has an application record.
// Inserts a new draft if none exists. Returns an error if something goes wrong.
func FirstFill(ctx context.Context, role string, userID int64, tx pgx.Tx) error {
	// Only run for students
	if role != "student" {
		log.Println("FirstFill skipped: user is not a student")
		return nil
	}

	query := `
    INSERT INTO applications (
        user_id,
        status,
        created_at,
        last_updated
    )
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (user_id) DO NOTHING;
    `

	ct, err := tx.Exec(ctx, query, userID, "Not Submitted", time.Now(), time.Now())
	if err != nil {
		log.Println("FirstFill: failed to insert application:", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		log.Printf("FirstFill: application already exists for user_id=%d\n", userID)
	} else {
		log.Printf("FirstFill: inserted new draft application for user_id=%d\n", userID)
	}

	return nil
}

func UpdateApplication(
	ctx context.Context,
	tx pgx.Tx,

	status string,
	dateOfBirth *time.Time,
	gender *string,
	homeDistrict *string,
	programme *string,
	registrationNumber *string,
	typeOfIntake *string,
	accommodation *string,
	parentGuardianStatus *string,
	guardianEmploymentStatus *string,
	householdMonthlyIncome *float64,
	otherFinancialSupport *string,
	reasonForBursary *string,
	submissionTimestamp *time.Time,
	userID int64,
) error {

	query := `
	UPDATE applications
	SET
		status                     = $1,
		date_of_birth              = $2,
		gender                     = $3,
		home_district              = $4,
		programme                  = $5,
		registration_number        = $6,
		type_of_intake             = $7,
		accommodation              = $8,
		parent_guardian_status     = $9,
		guardian_employment_status = $10,
		household_monthly_income   = $11,
		bursary_amount             = 0,
		other_financial_support    = $12,
		reason_for_bursary         = $13,
		submission_timestamp       = $14,
		last_updated               = NOW(),
		application_date           = $16
	WHERE user_id = $15
	  AND submission_timestamp IS NULL;
	`
	today := time.Now()
	ct, err := tx.Exec(
		ctx,
		query,
		status,
		dateOfBirth,
		gender,
		homeDistrict,
		programme,
		registrationNumber,
		typeOfIntake,
		accommodation,
		parentGuardianStatus,
		guardianEmploymentStatus,
		householdMonthlyIncome,
		otherFinancialSupport,
		reasonForBursary,
		submissionTimestamp,
		userID,
		time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()),
	)

	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("application already submitted or not found")
	}

	return nil
}
