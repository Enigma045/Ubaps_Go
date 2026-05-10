package Routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
	notifications "ubaps/Notifications"
	user_logs "ubaps/Audit_logs"
	"ubaps/services"
	"time"
	"ubaps/utils"
	middleware "ubaps/Middleware"
)

func Applicants(w http.ResponseWriter, r *http.Request) {
	var filters Handles.ReportFilters
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	
	err := json.NewDecoder(r.Body).Decode(&filters)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, limit, offset := utils.GetPaginationParams(r)
	year := r.URL.Query().Get("year")
	semester := r.URL.Query().Get("semester")
	month := r.URL.Query().Get("month")

	results, total, err := Handles.Applicants(Db.DB, ctx, filters, limit, offset, year, semester, month)
	if err != nil {
		log.Println("Error retriving applicants from database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(utils.PaginatedResponse{
		Data:  results,
		Total: total,
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		log.Println("Error encording applicants to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type activeStudent struct {
	StudentName string `json:"name"`
	StudentID   string `json:"id"`
}

func ConsiderStudent(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content", "application/json")
	ctx := r.Context()
	
	adminID, _ := middleware.UserIDFromContext(ctx)
	var adminIDPtr *int64
	if adminID != 0 {
		adminIDPtr = &adminID
	}

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var student activeStudent
	err = json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Println("Error decoding request body:", err)
		user_logs.Create_user_log(tx, adminIDPtr, "committee", "CONSIDER_STUDENT_FAILED", "Invalid JSON payload", "FAILED", time.Since(start), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := fmt.Sprintf("%s@unilia.ac.mw", student.StudentID)
	UserId, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println("Error getting user ID from email:", err)
		user_logs.Create_user_log(tx, adminIDPtr, "committee", "CONSIDER_STUDENT_FAILED", fmt.Sprintf("Student not found: %s", student.StudentID), "FAILED", time.Since(start), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Role, okay := middleware.RoleFromContext(ctx)
	if okay != true {
		log.Println("Error Retrieving Role from Context")
		user_logs.Create_user_log(tx, adminIDPtr, "committee", "CONSIDER_STUDENT_FAILED", "Authentication role error", "FAILED", time.Since(start), &UserId)
		http.Error(w, "Failed to retrieve role", http.StatusInternalServerError)
		return
	}

	results, err := Handles.ConsiderStudent(tx, ctx, UserId, Role)
	if err != nil {
		log.Println("Error Considering applicants from database", err)
		user_logs.Create_application_log(tx, adminIDPtr, Role, "CONSIDER_STUDENT_FAILED", err.Error(), "FAILED", time.Since(start), student.StudentID, nil, &UserId)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Audit Log (APPLICATION_LOG)
	user_logs.Create_application_log(tx, adminIDPtr, Role, fmt.Sprintf("%s_CONSIDERED_STUDENT", Role), fmt.Sprintf("Student:%s", student.StudentID), "SUCCESS", time.Since(start), student.StudentID, nil, &UserId)

	if err := tx.Commit(ctx); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Println("Error encoding applicants to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RejectStudent(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content", "application/json")
	ctx := r.Context()

	adminID, _ := middleware.UserIDFromContext(ctx)
	var adminIDPtr *int64
	if adminID != 0 {
		adminIDPtr = &adminID
	}

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var student activeStudent
	err = json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Println("Error decoding request body:", err)
		user_logs.Create_user_log(tx, adminIDPtr, "registrar", "BURSARY_REJECTION_FAILED", "Invalid JSON payload", "FAILED", time.Since(start), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := fmt.Sprintf("%s@unilia.ac.mw", student.StudentID)
	UserId, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println("Error getting user ID from email:", err)
		user_logs.Create_user_log(tx, adminIDPtr, "registrar", "BURSARY_REJECTION_FAILED", fmt.Sprintf("Student not found: %s", student.StudentID), "FAILED", time.Since(start), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Role, okay := middleware.RoleFromContext(ctx)
	if okay != true {
		log.Println("Error Retrieving Role from Context", err)
		user_logs.Create_user_log(tx, adminIDPtr, "registrar", "BURSARY_REJECTION_FAILED", "Authentication role error", "FAILED", time.Since(start), &UserId)
		http.Error(w, "Failed to retrieve role", http.StatusInternalServerError)
		return
	}

	results, err := Handles.RejectStudent(tx, ctx, UserId, Role)
	if err != nil {
		log.Println("Error rejecting applicants from database", err)
		user_logs.Create_application_log(tx, adminIDPtr, Role, "BURSARY_REJECTION_FAILED", err.Error(), "FAILED", time.Since(start), student.StudentID, nil, &UserId)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Final Rejection logic for registrar
	if Role == "registrar" {
		// 1. Notify Student (In-app)
		notifications.Send_notification(UserId, tx, "Your application for the bursary scheme was rejected at this time.", "Application Status")

		// 2. Notify Student (SMS)
		if phone, err := Handles.GetUserPhoneByID(UserId, tx); err == nil && phone != "" {
			smsMsg := "Your application for the bursary scheme was rejected at this time. Log in to your portal for more information."
			services.SendSMS(phone, smsMsg)
		}
		// Notify monitoring number (Intentional hardcoded value)
		services.SendSMS("0998111960", fmt.Sprintf("MONITOR: Student %s application was NOT selected.", student.StudentID))


		// 3. Email Administrator
		subject := "Bursary Rejection Finalized"
		body := fmt.Sprintf(`
			<h3>Bursary Rejection Alert</h3>
			<p>A student application has been officially rejected by the registrar.</p>
			<ul>
				<li><strong>Student Reg:</strong> %s</li>
				<li><strong>Status:</strong> Not Selected</li>
			</ul>
		`, student.StudentID)
		services.SendEmail("richardsambo94@gmail.com", subject, body)

		// 4. Audit Log (APPLICATION_LOG)
		user_logs.Create_application_log(tx, adminIDPtr, Role, "STUDENT_BURSARY_DISAPPROVED", fmt.Sprintf("Student:%s", student.StudentID), "SUCCESS", time.Since(start), student.StudentID, nil, &UserId)
	} else {
		// Intermediate disapproval log (APPLICATION_LOG)
		user_logs.Create_application_log(tx, adminIDPtr, Role, fmt.Sprintf("%s_DISAPPROVED_STUDENT", Role), fmt.Sprintf("Student:%s", student.StudentID), "SUCCESS", time.Since(start), student.StudentID, nil, &UserId)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Error committing transaction:", err)
		user_logs.Create_application_log(nil, adminIDPtr, Role, "BURSARY_REJECTION_COMMIT_FAILED", err.Error(), "FAILED", time.Since(start), student.StudentID, nil, &UserId)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Println("Error encoding applicants to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PayInstallment(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content", "application/json")
	ctx := r.Context()

	adminID, _ := middleware.UserIDFromContext(ctx)
	var adminIDPtr *int64
	if adminID != 0 {
		adminIDPtr = &adminID
	}

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var student activeStudent
	err = json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Println("Error decoding request body:", err)
		user_logs.Create_user_log(tx, adminIDPtr, "finance", "PAYMENT_FAILED", "Invalid JSON payload", "FAILED", time.Since(start), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := fmt.Sprintf("%s@unilia.ac.mw", student.StudentID)
	UserId, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println("Error getting user ID from email:", err)
		user_logs.Create_user_log(tx, adminIDPtr, "finance", "PAYMENT_FAILED", fmt.Sprintf("Student not found: %s", student.StudentID), "FAILED", time.Since(start), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := Handles.PayInstallment(tx, ctx, UserId)
	if err != nil {
		log.Println("Error processing payment in database", err)
		user_logs.Create_user_log(tx, adminIDPtr, "finance", "PAYMENT_FAILED", err.Error(), "FAILED", time.Since(start), &UserId)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 1. Notify Student (In-app)
	notifications.Send_notification(UserId, tx, "A payment installment has been successfully processed for your bursary account.", "Payment Processed")

	// 2. Notify Student (SMS)
	if phone, err := Handles.GetUserPhoneByID(UserId, tx); err == nil && phone != "" {
		smsMsg := "A payment installment has been successfully paid for your bursary account. Log in to your portal to view your updated statement."
		services.SendSMS("0998111960", smsMsg)
	}
	// Notify monitoring number (Intentional hardcoded value)
	//services.SendSMS("0998111960", fmt.Sprintf("MONITOR: Payment processed for student %s.", student.StudentID))


	// 3. Get Amount for Audit Log
	var amountStr string
	err = tx.QueryRow(ctx, "SELECT bursary_amount FROM applications WHERE user_id = $1", UserId).Scan(&amountStr)
	if err != nil {
		log.Println("Error fetching amount for log:", err)
	}
	var amount float64
	fmt.Sscanf(amountStr, "%f", &amount)

	// 3. Audit Log (PAYMENT_LOG)
	user_logs.Create_payment_log(tx, adminIDPtr, "finance", "BURSARY_PAYMENT_PROCESSED", fmt.Sprintf("Student:%s, Amount:MWK %s", student.StudentID, amountStr), "SUCCESS", time.Since(start), amount, student.StudentID, &UserId)

	// 4. Email Administrator
	subject := "Bursary Payment Processed"
	body := fmt.Sprintf(`
		<h3>Installment Payment Alert</h3>
		<p>A bursary installment has been successfully processed for a student.</p>
		<ul>
			<li><strong>Student Reg:</strong> %s</li>
			<li><strong>Amount:</strong> MWK %s</li>
			<li><strong>Status:</strong> Paid</li>
		</ul>
	`, student.StudentID, amountStr)
	services.SendEmail("richardsambo94@gmail.com", subject, body)

	if err := tx.Commit(ctx); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Println("Error encoding payment response to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func RollbackSelection(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	adminID, _ := middleware.UserIDFromContext(ctx)
	Role, _ := middleware.RoleFromContext(ctx)

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var student activeStudent
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := fmt.Sprintf("%s@unilia.ac.mw", student.StudentID)
	UserId, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	msg, err := Handles.RollbackSelection(tx, ctx, UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 1. Audit Logs
	user_logs.Create_user_log(tx, &adminID, Role, "ROLLBACK_SELECTION", fmt.Sprintf("Student:%s", student.StudentID), "SUCCESS", time.Since(start), &UserId)
	user_logs.Create_application_log(tx, &adminID, Role, "ROLLBACK_SELECTION", fmt.Sprintf("Student:%s", student.StudentID), "SUCCESS", time.Since(start), student.StudentID, nil, &UserId)

	// 2. Notifications
	// To Initiator (Registrar)
	notifications.Send_notification(adminID, tx, fmt.Sprintf("You have successfully rolled back the selection for student %s.", student.StudentID), "Rollback Successful")

	// To Deans and Finance Office
	staffRoles := []string{"dean_of_student", "dean_of_facult", "dean_of_science", "finance_office"}
	for _, role := range staffRoles {
		if staffIDs, err := Handles.GetUserIDsOfDifferentTypes(tx, role); err == nil {
			notifications.BroadcastNotification(staffIDs, tx, fmt.Sprintf("The selection for student %s has been rolled back and funds restored.", student.StudentID), "Selection Rollback Alert")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

type commentRequest struct {
	StudentID string `json:"id"`
	Comment   string `json:"comment"`
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	adminID, _ := middleware.UserIDFromContext(ctx)
	Role, _ := middleware.RoleFromContext(ctx)

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var req commentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := fmt.Sprintf("%s@unilia.ac.mw", req.StudentID)
	UserId, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Get commenter name
	commenterName, _ := Handles.GetEmailByUserID(adminID, tx)

	err = Handles.AddComment(tx, ctx, UserId, req.Comment, commenterName, Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 1. Audit Logs
	user_logs.Create_user_log(tx, &adminID, Role, "APPLICATION_COMMENT_ADDED", fmt.Sprintf("Student:%s, Comment:%s", req.StudentID, req.Comment), "SUCCESS", time.Since(start), &UserId)
	user_logs.Create_application_log(tx, &adminID, Role, "APPLICATION_COMMENT_ADDED", fmt.Sprintf("Student:%s, Comment:%s", req.StudentID, req.Comment), "SUCCESS", time.Since(start), req.StudentID, nil, &UserId)

	// 2. Notifications to other staff members
	allStaffRoles := []string{"registrar", "dean_of_student", "dean_of_facult", "dean_of_science", "finance_office"}
	for _, role := range allStaffRoles {
		if staffIDs, err := Handles.GetUserIDsOfDifferentTypes(tx, role); err == nil {
			// Filter out the commenter themselves
			var filteredIDs []int64
			for _, id := range staffIDs {
				if id != adminID {
					filteredIDs = append(filteredIDs, id)
				}
			}
			if len(filteredIDs) > 0 {
				notifications.BroadcastNotification(filteredIDs, tx, fmt.Sprintf("%s (%s) made a comment on student %s's application.", commenterName, Role, req.StudentID), "New Application Comment")
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully"})
}
