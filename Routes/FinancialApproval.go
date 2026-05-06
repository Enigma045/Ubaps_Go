package Routes

import (
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
	middleware "ubaps/Middleware"
	"ubaps/utils"
	user_logs "ubaps/Audit_logs"
	notifications "ubaps/Notifications"
	"time"
	"fmt"
)	

func Approval(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx := r.Context()
	tx, err := Db.DB.Begin(ctx)

	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Database transaction failed", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	formdata, err := utils.Formdata(r)
	if err != nil {
		log.Println("Failed to retrieve formdata:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	userid, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		log.Println("Failed to retrieve userid from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Assuming the request amount and details are in formdata
	amountStr, _ := utils.GetFormValue(formdata, "amount")
	details, _ := utils.GetFormValue(formdata, "detail")

	err = Handles.Db_operation(tx, ctx, userid, 0, formdata)
	if err != nil {
		log.Println("Database operation failed:", err)
		user_logs.Create_user_log(tx, &userid, "student", "FINANCIAL_REQUEST_FAILED", err.Error(), "FAILED", time.Since(start), &userid)
		http.Error(w, "Failed to submit financial request", http.StatusInternalServerError)
		return
	}

	// Audit Log
	user_logs.Create_user_log(tx, &userid, "student", "FINANCIAL_REQUEST_SUBMITTED", fmt.Sprintf("Amount:%s, Details:%s", amountStr, details), "SUCCESS", time.Since(start), &userid)

	// Notification to Student
	notifications.Send_notification(userid, tx, "Your financial request has been successfully submitted and is awaiting approval.", "Financial Request Submitted")

	// Broadcast to Finance Officers
	email, _ := Handles.GetEmailByUserID(userid, tx)
	if financeIDs, err := Handles.GetUserIDsOfDifferentTypes(tx, "student"); err == nil {
		// Note: Filter for finance_office if possible, otherwise broadcast to all staff
		notifications.BroadcastNotification(financeIDs, tx, fmt.Sprintf("A new financial request has been submitted by student %s.", email), "New Financial Request")
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("Failed to commit transaction:", err)
		http.Error(w, "Database commit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Fees statement successfully sent"))
}
