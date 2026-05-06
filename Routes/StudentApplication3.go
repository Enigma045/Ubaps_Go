package Routes

import (
	"log"
	"net/http"
	"time"
	"ubaps/Db"
	"ubaps/Handles"
	middleware "ubaps/Middleware"
	"ubaps/utils"
	user_logs "ubaps/Audit_logs"
	notifications "ubaps/Notifications"
	"fmt"
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
