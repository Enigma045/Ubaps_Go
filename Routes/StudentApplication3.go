package Routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	user_logs "ubaps/Audit_logs"
	"ubaps/Db"
	"ubaps/Handles"
	middleware "ubaps/Middleware"
	notifications "ubaps/Notifications"
	"ubaps/utils"

	"github.com/jackc/pgx/v5"
)

func SubmitForm(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	Ctx := r.Context()

	tx, err := Db.DB.Begin(Ctx)
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(Ctx)

	userId, ok := middleware.UserIDFromContext(Ctx)
	if !ok {
		log.Println("User not authenticated")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	email, err := Handles.GetEmailByUserID(userId, tx)
	if err != nil {
		log.Println("Error getting email:", err)
		user_logs.Create_application_log(tx, &userId, "student", "APPLICATION_SUBMISSION_FAILED", "Email lookup failure", "FAILED", time.Since(start), "UNKNOWN", nil, &userId)
		http.Error(w, "Failed to get email", http.StatusInternalServerError)
		return
	}

	RegNumber := Handles.GetRegNumberFromEmail(email)

	dob, err := Handles.TimePtr(r.FormValue("dob"))
	if err != nil {
		log.Println("Invalid date of birth:", err)
		user_logs.Create_application_log(tx, &userId, "student", "APPLICATION_SUBMISSION_FAILED", "Invalid DOB", "FAILED", time.Since(start), RegNumber, nil, &userId)
		http.Error(w, "Invalid date of birth", http.StatusBadRequest)
		return
	}

	submission := time.Now()

	err = utils.UpdateApplication(
		Ctx,
		tx,
		"submitted",
		dob,
		Handles.StrPtr(r.FormValue("gender")),
		Handles.StrPtr(r.FormValue("HomeDistrict")),
		Handles.StrPtr("Computer Engineering"),
		Handles.StrPtr(RegNumber),
		Handles.StrPtr(r.FormValue("Accomodation")),
		Handles.StrPtr(r.FormValue("Gurdian Status")),
		Handles.StrPtr(r.FormValue("Guardian Employment Status")),
		Handles.StrPtr(r.FormValue("otherSupport")),
		Handles.StrPtr(r.FormValue("Reason")),
		Handles.StrPtr(r.FormValue("feeResponsibility")),
		Handles.StrPtr(r.FormValue("financialHardship")),
		Handles.StrPtr(r.FormValue("impactOfNoBursary")),
		&submission,
		userId,
	)
	if err != nil {
		log.Println("UpdateApplication Error:", err)
		user_logs.Create_application_log(tx, &userId, "student", "APPLICATION_SUBMISSION_FAILED", err.Error(), "FAILED", time.Since(start), RegNumber, nil, &userId)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 1. Audit Log (Success)
	user_logs.Create_application_log(tx, &userId, "student", "BURSARY_FORM_SUBMITTED", fmt.Sprintf("Reg:%s", RegNumber), "SUCCESS", time.Since(start), RegNumber, nil, &userId)

	// 2. Notify Student
	notifications.Send_notification(userId, tx, "You have successfully submitted your bursary application form. You will be notified of the outcome soon.", "Application Submitted")

	// 3. Notify Staff (Excluding students and admins)
	if staffIDs, err := Handles.GetUserIDsOfDifferentTypes(tx, "student"); err == nil {
		notifications.BroadcastNotification(staffIDs, tx, fmt.Sprintf("A new bursary application has been submitted by student %s (%s).", email, RegNumber), "New Application Received")
	}

	if err := tx.Commit(Ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You have successfully submitted the application form"))
}

func GetMyApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var data struct {
		Status                   string    `json:"status"`
		Dob                      *time.Time `json:"dob"`
		Gender                   string    `json:"gender"`
		HomeDistrict             string    `json:"home_district"`
		Accommodation            string    `json:"accommodation"`
		GuardianStatus           string    `json:"guardian_status"`
		GuardianEmploymentStatus string    `json:"guardian_employment_status"`
		OtherSupport             string    `json:"other_support"`
		Reason                   string    `json:"reason"`
		FeeResponsibility        string    `json:"fee_responsibility"`
		FinancialHardship        string    `json:"financial_hardship"`
		ImpactOfNoBursary        string    `json:"impact_of_no_bursary"`
	}

	query := `
		SELECT 
			COALESCE(status, 'not submitted'), date_of_birth, COALESCE(gender, ''), 
			COALESCE(home_district, ''), COALESCE(accommodation, ''), 
			COALESCE(parent_guardian_status, ''), COALESCE(guardian_employment_status, ''), 
			COALESCE(other_financial_support, ''), COALESCE(reason_for_bursary, ''), 
			COALESCE(fee_responsibility, ''), COALESCE(financial_hardship, ''), 
			COALESCE(impact_of_no_bursary, '')
		FROM applications
		WHERE user_id = $1
	`
	
	err := Db.DB.QueryRow(ctx, query, userID).Scan(
		&data.Status, &data.Dob, &data.Gender, &data.HomeDistrict,
		&data.Accommodation, &data.GuardianStatus, &data.GuardianEmploymentStatus,
		&data.OtherSupport, &data.Reason, &data.FeeResponsibility,
		&data.FinancialHardship, &data.ImpactOfNoBursary,
	)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("GetMyApplication Error:", err)
		http.Error(w, "Failed to fetch application", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func GetApplicantHardship(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Only authorized staff should see this
	_, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	regNumber := r.URL.Query().Get("id")
	if regNumber == "" {
		http.Error(w, "Missing registration number", http.StatusBadRequest)
		return
	}

	var data struct {
		Reason            string `json:"reason"`
		FeeResponsibility string `json:"fee_responsibility"`
		FinancialHardship string `json:"financial_hardship"`
		ImpactOfNoBursary string `json:"impact_of_no_bursary"`
	}

	query := `
		SELECT 
			COALESCE(reason_for_bursary, ''), 
			COALESCE(fee_responsibility, ''), 
			COALESCE(financial_hardship, ''), 
			COALESCE(impact_of_no_bursary, '')
		FROM applications
		WHERE registration_number = $1
	`
	
	err := Db.DB.QueryRow(ctx, query, regNumber).Scan(
		&data.Reason, &data.FeeResponsibility,
		&data.FinancialHardship, &data.ImpactOfNoBursary,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Applicant not found", http.StatusNotFound)
		} else {
			log.Println("GetApplicantHardship Error:", err)
			http.Error(w, "Failed to fetch hardship details", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
