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
	var Pplicants []string
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	
	err := json.NewDecoder(r.Body).Decode(&Pplicants)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, limit, offset := utils.GetPaginationParams(r)

	results, total, err := Handles.Applicants(Db.DB, ctx, Pplicants, limit, offset)
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
		user_logs.Create_user_log(tx, nil, "committee", "CONSIDER_STUDENT_FAILED", "Invalid JSON payload", "FAILED", time.Since(start))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := fmt.Sprintf("%s@unilia.ac.mw", student.StudentID)
	UserId, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println("Error getting user ID from email:", err)
		user_logs.Create_user_log(tx, nil, "committee", "CONSIDER_STUDENT_FAILED", fmt.Sprintf("Student not found: %s", student.StudentID), "FAILED", time.Since(start))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Role, okay := middleware.RoleFromContext(ctx)
	if okay != true {
		log.Println("Error Retrieving Role from Context")
		user_logs.Create_user_log(tx, &UserId, "committee", "CONSIDER_STUDENT_FAILED", "Authentication role error", "FAILED", time.Since(start))
		http.Error(w, "Failed to retrieve role", http.StatusInternalServerError)
		return
	}

	results, err := Handles.ConsiderStudent(tx, ctx, UserId, Role)
	if err != nil {
		log.Println("Error Considering applicants from database", err)
		user_logs.Create_user_log(tx, &UserId, Role, "CONSIDER_STUDENT_FAILED", err.Error(), "FAILED", time.Since(start))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Audit Log
	user_logs.Create_user_log(tx, &UserId, Role, fmt.Sprintf("%s_CONSIDERED_STUDENT", Role), fmt.Sprintf("Student:%s", student.StudentID), "SUCCESS", time.Since(start))

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
		user_logs.Create_user_log(tx, nil, "registrar", "BURSARY_REJECTION_FAILED", "Invalid JSON payload", "FAILED", time.Since(start))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := fmt.Sprintf("%s@unilia.ac.mw", student.StudentID)
	UserId, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println("Error getting user ID from email:", err)
		user_logs.Create_user_log(tx, nil, "registrar", "BURSARY_REJECTION_FAILED", fmt.Sprintf("Student not found: %s", student.StudentID), "FAILED", time.Since(start))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Role, okay := middleware.RoleFromContext(ctx)
	if okay != true {
		log.Println("Error Retrieving Role from Context", err)
		user_logs.Create_user_log(tx, &UserId, "registrar", "BURSARY_REJECTION_FAILED", "Authentication role error", "FAILED", time.Since(start))
		http.Error(w, "Failed to retrieve role", http.StatusInternalServerError)
		return
	}

	results, err := Handles.RejectStudent(tx, ctx, UserId, Role)
	if err != nil {
		log.Println("Error rejecting applicants from database", err)
		user_logs.Create_user_log(tx, &UserId, Role, "BURSARY_REJECTION_FAILED", err.Error(), "FAILED", time.Since(start))
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
			services.SendSMS("0998111960", smsMsg)
		}

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

		// 4. Audit Log
		user_logs.Create_user_log(tx, &UserId, Role, "STUDENT_BURSARY_REJECTION", fmt.Sprintf("Student:%s", student.StudentID), "SUCCESS", time.Since(start))
	} else {
		// Intermediate rejection log
		user_logs.Create_user_log(tx, &UserId, Role, fmt.Sprintf("%s_REJECTED_STUDENT", Role), fmt.Sprintf("Student:%s", student.StudentID), "SUCCESS", time.Since(start))
	}

	tx.Commit(ctx)
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
		user_logs.Create_user_log(tx, nil, "finance", "PAYMENT_FAILED", "Invalid JSON payload", "FAILED", time.Since(start))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := fmt.Sprintf("%s@unilia.ac.mw", student.StudentID)
	UserId, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println("Error getting user ID from email:", err)
		user_logs.Create_user_log(tx, nil, "finance", "PAYMENT_FAILED", fmt.Sprintf("Student not found: %s", student.StudentID), "FAILED", time.Since(start))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := Handles.PayInstallment(tx, ctx, UserId)
	if err != nil {
		log.Println("Error processing payment in database", err)
		user_logs.Create_user_log(tx, &UserId, "finance", "PAYMENT_FAILED", err.Error(), "FAILED", time.Since(start))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 1. Notify Student (In-app)
	notifications.Send_notification(UserId, tx, "A payment installment has been successfully processed for your bursary account.", "Payment Processed")

	// 2. Notify Student (SMS)
	if phone, err := Handles.GetUserPhoneByID(UserId, tx); err == nil && phone != "" {
		smsMsg := "A payment installment has been successfully processed for your bursary account. Log in to your portal to view your updated statement."
		services.SendSMS("0998111960", smsMsg)
	}

	// 3. Audit Log
	user_logs.Create_user_log(tx, &UserId, "finance", "BURSARY_PAYMENT_PROCESSED", fmt.Sprintf("Student:%s", student.StudentID), "SUCCESS", time.Since(start))

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
