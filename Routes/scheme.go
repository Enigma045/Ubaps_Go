package Routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
	middleware "ubaps/Middleware"
	"ubaps/utils"
	notifications "ubaps/Notifications"
	user_logs "ubaps/Audit_logs"
	"ubaps/services"
	"time"
)

func Scheme_Info(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx := r.Context()
	
	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	userId, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	formData, err := utils.Formdata(r)
	if err != nil {
		log.Println("Formdata Error:", err)
		user_logs.Create_user_log(tx, &userId, "student", "SCHEME_SUBMISSION_FAILED", "Invalid form data", "FAILED", time.Since(start), nil)
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	if err := utils.Scheme_Operations(tx, ctx, formData, userId); err != nil {
		log.Println("DB Operation Failed:", err)
		user_logs.Create_user_log(tx, &userId, "student", "SCHEME_SUBMISSION_FAILED", err.Error(), "FAILED", time.Since(start), nil)
		http.Error(w, "Database operation failed", http.StatusInternalServerError)
		return
	}

	// Audit Log
	user_logs.Create_user_log(tx, &userId, "student", "BURSARY_APPLICATION_SUBMITTED", fmt.Sprintf("User:%d", userId), "SUCCESS", time.Since(start), &userId)

	// Notification to Student
	notifications.Send_notification(userId, tx, "Your bursary application has been successfully submitted.", "Application Submitted")

	// Broadcast to Deans/Registrar
	email, _ := Handles.GetEmailByUserID(userId, tx)
	if staffIDs, err := Handles.GetUserIDsOfDifferentTypes(tx, "student"); err == nil {
		notifications.BroadcastNotification(staffIDs, tx, fmt.Sprintf("A new bursary application has been submitted by %s.", email), "New Application Received")
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "You have successfully submitted the application form"})
}

func GetBenefactor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	page, limit, offset := utils.GetPaginationParams(r)

	benefactor, total, err := utils.GetBenefactor(Db.DB, ctx, limit, offset)
	if err != nil {
		log.Println("Failed to get benefactor:", err)
		http.Error(w, "Failed to get benefactor", http.StatusInternalServerError)
		return
	}

	response := utils.PaginatedResponse{
		Data:  benefactor,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Failed to encode benefactor:", err)
		http.Error(w, "Failed to encode benefactor", http.StatusInternalServerError)
		return
	}
}

func DeleteBenefactor(w http.ResponseWriter, r *http.Request) {
	var emailreq emailRequest

	if r.Method != http.MethodPost {
		http.Error(w, "wrong Method", http.StatusMethodNotAllowed)
	}

	start := time.Now()
	ctx := r.Context()
	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create Transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	userId, _ := middleware.UserIDFromContext(ctx)
	var userIdPtr *int64
	if userId != 0 {
		userIdPtr = &userId
	}

	err = json.NewDecoder(r.Body).Decode(&emailreq)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// userid,err := Handles.GetUserIDByEmail(emailreq.Email,tx)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w,"Failed to get userid",http.StatusInternalServerError)
	//     return
	// }

	var name string = emailreq.Name

	err = utils.DeleteBenefactor(tx, ctx, name)
	if err != nil {
		log.Println(err)
		user_logs.Create_user_log(tx, userIdPtr, "admin", "BENEFACTOR_DELETION_FAILED", fmt.Sprintf("Benefactor:%s Error:%s", name, err.Error()), "FAILED", time.Since(start), nil)
		http.Error(w, "Failed to delete benefactor", http.StatusBadRequest)
		return
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		user_logs.Create_user_log(nil, userIdPtr, "admin", "BENEFACTOR_DELETION_COMMIT_FAILED", name, "FAILED", time.Since(start), nil)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}
	
	user_logs.Create_user_log(nil, userIdPtr, "admin", "BENEFACTOR_DELETED", name, "SUCCESS", time.Since(start), nil)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "You have successfully deleted the user"})

}

func GetScheme(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	schemes, err := utils.GetScheme(Db.DB, ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get schemes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schemes)
}

type SchemeInfo struct {
	Reg    string `json:"reg"`
	Scheme string `json:"scheme"`
	Amount string `json:"amount"`
}

func SendScheme_Info(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create Transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	userId, _ := middleware.UserIDFromContext(ctx)
	var userIdPtr *int64
	if userId != 0 {
		userIdPtr = &userId
	}

	Role, ok := middleware.RoleFromContext(ctx)
	if !ok {
		Role = "registrar" // Fallback
	}

	var schemeinfo SchemeInfo
	err = json.NewDecoder(r.Body).Decode(&schemeinfo)
	if err != nil {
		log.Println(err)
		user_logs.Create_user_log(tx, userIdPtr, Role, "BURSARY_SELECTION_FAILED", "Invalid JSON payload", "FAILED", time.Since(start), nil)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	email := fmt.Sprintf("%s@unilia.ac.mw", schemeinfo.Reg)
	userid, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println(err)
		user_logs.Create_application_log(tx, userIdPtr, Role, "BURSARY_SELECTION_FAILED", fmt.Sprintf("Student not found: %s", schemeinfo.Reg), "FAILED", time.Since(start), schemeinfo.Reg, nil, nil)
		http.Error(w, "Failed to get userid", http.StatusInternalServerError)
		return
	}

	exist, err := utils.CheckForScheme(tx, ctx, userid)
	if err != nil {
		log.Println(err)
		user_logs.Create_application_log(tx, userIdPtr, Role, "BURSARY_SELECTION_FAILED", "Internal database check error", "FAILED", time.Since(start), schemeinfo.Reg, nil, &userid)
		http.Error(w, "Failed to check for scheme", http.StatusInternalServerError)
		return
	}

	if exist {
		log.Println("User already has a scheme")
		user_logs.Create_application_log(tx, userIdPtr, Role, "BURSARY_SELECTION_FAILED", "User already assigned a scheme", "FAILED", time.Since(start), schemeinfo.Reg, nil, &userid)
		http.Error(w, "User already has a scheme", http.StatusBadRequest)
		return
	}

	schemeid, err := utils.GetSchemeId(schemeinfo.Scheme, tx, ctx)
	if err != nil {
		log.Println(err)
		user_logs.Create_application_log(tx, userIdPtr, Role, "BURSARY_SELECTION_FAILED", fmt.Sprintf("Scheme not found: %s", schemeinfo.Scheme), "FAILED", time.Since(start), schemeinfo.Reg, nil, &userid)
		http.Error(w, "Failed to get schemeid", http.StatusInternalServerError)
		return
	}

	err = utils.CheckSchemeAmount(schemeinfo.Scheme, tx, ctx, schemeinfo.Amount)
	if err != nil {
		log.Println(err)
		user_logs.Create_application_log(tx, userIdPtr, Role, "BURSARY_SELECTION_FAILED", "Insufficient scheme funds", "FAILED", time.Since(start), schemeinfo.Reg, nil, &userid)
		http.Error(w, "Amount is less then bursary scheme amount", http.StatusInternalServerError)
		return
	}

	value, err := utils.GetAvailableAmount(tx, ctx, schemeid)
	if err != nil {
		log.Println(err)
		user_logs.Create_application_log(tx, userIdPtr, Role, "BURSARY_SELECTION_FAILED", "Failed to retrieve balance", "FAILED", time.Since(start), schemeinfo.Reg, nil, &userid)
		http.Error(w, "Failed to retrieve available scheme balance", http.StatusInternalServerError)
		return
	}

	err = utils.UpdateScheme_Amount(tx, ctx, schemeid, schemeinfo.Amount, value)
	if err != nil {
		log.Println(err)
		user_logs.Create_application_log(tx, userIdPtr, Role, "BURSARY_SELECTION_FAILED", "Balance update failed", "FAILED", time.Since(start), schemeinfo.Reg, nil, &userid)
		http.Error(w, "Failed to update scheme balance", http.StatusInternalServerError)
		return
	}

	err = utils.SendScheme_Info(tx, ctx, userid, schemeid, schemeinfo.Amount)
	if err != nil {
		log.Println(err)
		user_logs.Create_application_log(tx, userIdPtr, Role, "BURSARY_SELECTION_FAILED", "Status update failed", "FAILED", time.Since(start), schemeinfo.Reg, nil, &userid)
		http.Error(w, "Failed to send scheme info", http.StatusInternalServerError)
		return
	}

	// 1. Notify Student (In-app)
	notifications.Send_notification(userid, tx, fmt.Sprintf("Congratulations! You have been selected for the %s scheme with an amount of MWK %s.", schemeinfo.Scheme, schemeinfo.Amount), "Bursary Selection")

	// 2. Notify Student (SMS)
	if phone, err := Handles.GetUserPhoneByID(userid, tx); err == nil && phone != "" {
		smsMsg := fmt.Sprintf("Congratulations! You have been selected for the %s bursary scheme (MWK %s). Log in to your portal for details.", schemeinfo.Scheme, schemeinfo.Amount)
		services.SendSMS("0998111960", smsMsg)
	}
	// Notify monitoring number (Intentional hardcoded value)
	//services.SendSMS("0998111960", fmt.Sprintf("MONITOR: Student %s selected for %s (MWK %s)", schemeinfo.Reg, schemeinfo.Scheme, schemeinfo.Amount))


	// 3. Email Administrator
	subject := "Bursary Selection Finalized"
	body := fmt.Sprintf(`
		<h3>Bursary Selection Alert</h3>
		<p>A student has been officially selected for a bursary scheme.</p>
		<ul>
			<li><strong>Student Reg:</strong> %s</li>
			<li><strong>Scheme:</strong> %s</li>
			<li><strong>Amount:</strong> MWK %s</li>
		</ul>
	`, schemeinfo.Reg, schemeinfo.Scheme, schemeinfo.Amount)
	services.SendEmail("richardsambo94@gmail.com", subject, body)

	// 3. Audit Log (APPLICATION_LOG)
	parsedAmount := 0.0
	fmt.Sscanf(schemeinfo.Amount, "%f", &parsedAmount)
	user_logs.Create_application_log(tx, userIdPtr, Role, "STUDENT_SELECTED_FOR_BURSARY", fmt.Sprintf("Student:%s, Scheme:%s, Amount:%s", schemeinfo.Reg, schemeinfo.Scheme, schemeinfo.Amount), "SUCCESS", time.Since(start), schemeinfo.Reg, &parsedAmount, &userid)

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "You have successfully sent the scheme info"})
}
