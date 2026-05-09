package Routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"ubaps/Db"
	"ubaps/Handles"
	"ubaps/utils"
	user_logs "ubaps/Audit_logs"
	middleware "ubaps/Middleware"
	notifications "ubaps/Notifications"
)

func Fees(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	performerID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx := r.Context()
	tx, err := Db.DB.Begin(ctx)

	if err != nil {
		http.Error(w, "Database Failed transactions Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	reg := r.FormValue("student_id")
	email := fmt.Sprintf("%s@unilia.ac.mw", reg)

	data, err := utils.Formdata(r)
	if err != nil {
		http.Error(w, "Formdata Error", http.StatusInternalServerError)
		return
	}
	log.Println(email)
	student_id, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get student_id", http.StatusInternalServerError)
		return
	}
	err = utils.Finance_Operations(tx, ctx, data, student_id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to perform finance operations", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		user_logs.Create_user_log(nil, &performerID, "finance_office", "DOSSIER_RESPONSE_FAILED", fmt.Sprintf("Failed to submit fees statement for %s: %s", reg, err.Error()), "FAILED", time.Since(start), &student_id)
		http.Error(w, "Database Failed to commit Error", http.StatusInternalServerError)
		return
	}
	
	user_logs.Create_user_log(nil, &performerID, "finance_office", "DOSSIER_RESPONSE_SUCCESS", fmt.Sprintf("Fees statement successfully sent for %s", reg), "SUCCESS", time.Since(start), &student_id)
	w.Write([]byte("Fees statement successfully sent"))
}

func Request_Statement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx := r.Context()

	// Parse student_id (registration number) from request body (JSON)
	var payload struct {
		Reg string `json:"reg"`
	}
	if err := utils.DecodeJSON(r, &payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.Reg == "" {
		http.Error(w, "Student registration number is required", http.StatusBadRequest)
		return
	}

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	email := fmt.Sprintf("%s@unilia.ac.mw", payload.Reg)
	studentID, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println("Error resolving student ID:", err)
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	if err := utils.Request_Finance_Statement(tx, ctx, studentID); err != nil {
		log.Println("Error requesting statement:", err)
		http.Error(w, "Failed to request statement", http.StatusInternalServerError)
		return
	}

	// Notify Finance Office
	financeUserIDs, err := Handles.GetUserIDsByRole(tx, "finance_office")
	if err == nil && len(financeUserIDs) > 0 {
		performerRole, _ := middleware.RoleFromContext(r.Context())
		notifications.BroadcastNotification(
			financeUserIDs, 
			tx, 
			fmt.Sprintf("A new financial statement has been requested for %s by %s", payload.Reg, performerRole), 
			"Statement Requested",
		)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Statement request sent successfully"}`))
}

func GetFinancialHistory(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var payload struct {
		Reg string `json:"reg"`
	}
	if err := utils.DecodeJSON(r, &payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.Reg == "" {
		http.Error(w, "Student registration number is required", http.StatusBadRequest)
		return
	}

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	email := fmt.Sprintf("%s@unilia.ac.mw", payload.Reg)
	studentID, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println("Error resolving student ID:", err)
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	query := `
		SELECT 
			semester, 
			payment_date, 
			details, 
			payment_amount, 
			full_installment, 
			request, 
			updated_at 
		FROM financial_history 
		WHERE student_id = $1 
		ORDER BY updated_at DESC
	`

	rows, err := tx.Query(ctx, query, studentID)
	if err != nil {
		log.Println("Error fetching financial history:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type HistoryItem struct {
		Semester        string    `json:"semester"`
		Date            time.Time `json:"date"`
		Details         string    `json:"details"`
		Amount          float64   `json:"amount"`
		FullInstallment float64   `json:"full_installment"`
		Status          string    `json:"status"`
		UpdatedAt       time.Time `json:"updated_at"`
	}

	var history []HistoryItem
	for rows.Next() {
		var item HistoryItem
		if err := rows.Scan(
			&item.Semester,
			&item.Date,
			&item.Details,
			&item.Amount,
			&item.FullInstallment,
			&item.Status,
			&item.UpdatedAt,
		); err != nil {
			log.Println("Error scanning history row:", err)
			continue
		}
		history = append(history, item)
	}

	performerID, _ := middleware.UserIDFromContext(r.Context())
	user_logs.Create_user_log(nil, &performerID, "admin", "DOSSIER_REQUEST_SUCCESS", fmt.Sprintf("Viewed financial dossier for %s", payload.Reg), "SUCCESS", time.Since(start), &studentID)

	if err := tx.Commit(ctx); err != nil {
		log.Println("Error committing transaction:", err)
	}

	json.NewEncoder(w).Encode(history)
}
