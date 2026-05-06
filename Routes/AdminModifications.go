package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
	middleware "ubaps/Middleware"
	"ubaps/utils"
	user_logs "ubaps/Audit_logs"
	notifications "ubaps/Notifications"
	"ubaps/services"
	"time"
	"fmt"
	"crypto/rand"
	"encoding/hex"
)

type emailRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Getuserdetails(w http.ResponseWriter, r *http.Request) {
	// Allow only GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Failed to get user id", http.StatusUnauthorized)
		return
	}

	page, limit, offset := utils.GetPaginationParams(r)

	details, total, err := utils.ReciveDetails(Db.DB, ctx, userID, limit, offset)
	if err != nil {
		log.Println("Error fetching userDetails:", err)
		http.Error(w, "Failed to fetch userDetails", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.PaginatedResponse{
		Data:  details,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var emailreq emailRequest

	if r.Method != http.MethodPost {
		http.Error(w, "wrong Method", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create Transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	adminID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&emailreq)
	if err != nil {
		log.Println(err)
		user_logs.Create_user_log(tx, &adminID, "admin", "DELETE_ACCOUNT_FAILED", "Invalid JSON", "FAILED", time.Since(start), nil)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userid, err := Handles.GetUserIDByEmail(emailreq.Email, tx)
	if err != nil {
		log.Println(err)
		user_logs.Create_user_log(tx, &adminID, "admin", "DELETE_ACCOUNT_FAILED", fmt.Sprintf("User not found: %s", emailreq.Email), "FAILED", time.Since(start), nil)
		http.Error(w, "Failed to get userid", http.StatusInternalServerError)
		return
	}

	err = utils.DeleteUser(tx, ctx, userid)
	if err != nil {
		log.Println(err)
		user_logs.Create_user_log(tx, &adminID, "admin", "DELETE_ACCOUNT_FAILED", err.Error(), "FAILED", time.Since(start), &userid)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Audit Log
	user_logs.Create_user_log(tx, &adminID, "admin", "USER_ACCOUNT_DELETED", fmt.Sprintf("Deleted User:%s", emailreq.Email), "SUCCESS", time.Since(start), &userid)

	// Notify all admins about the deletion
	if adminIDs, err := Handles.GetAdmins(tx); err == nil {
		notifications.BroadcastNotification(adminIDs, tx, fmt.Sprintf("Admin has deleted the account of user: %s", emailreq.Email), "Account Deleted")
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You have successfully deleted the user"))
}

func AdminUpdateUserAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()
	ctx := r.Context()
	adminID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		OriginalEmail string `json:"originalEmail"`
		First         string `json:"first"`
		Last          string `json:"last"`
		Email         string `json:"email"`
		Phone         string `json:"phone"`
		Status        string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	targetUserID, err := Handles.GetUserIDByEmail(req.OriginalEmail, tx)
	if err != nil {
		http.Error(w, "Target user not found", http.StatusNotFound)
		return
	}

	err = Handles.AdminUpdateUser(ctx, tx, targetUserID, req.First, req.Last, req.Email, req.Phone, req.Status)
	if err != nil {
		log.Println("Error in AdminUpdateUser:", err)
		user_logs.Create_user_log(tx, &adminID, "admin", "ADMIN_USER_UPDATE_FAILED", err.Error(), "FAILED", time.Since(start), &targetUserID)
		tx.Commit(ctx)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Audit Log
	user_logs.Create_user_log(tx, &adminID, "admin", "ADMIN_USER_UPDATE_SUCCESS", fmt.Sprintf("Updated user %s: Status=%s", req.Email, req.Status), "SUCCESS", time.Since(start), &targetUserID)

	// Notifications
	notifications.Send_notification(adminID, tx, fmt.Sprintf("You have successfully updated user: %s", req.Email), "User Updated")
	notifications.Send_notification(targetUserID, tx, "Your account settings have been updated by an administrator.", "Account Updated")

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func TriggerPasswordResetAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()
	ctx := r.Context()
	adminID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	targetUserID, err := Handles.GetUserIDByEmail(req.Email, tx)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Generate reset token (reusing logic pattern)
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)
	expiresAt := time.Now().Add(1 * time.Hour)

	EnsurePasswordResetsTable(ctx)
	_, err = tx.Exec(ctx, `
		INSERT INTO password_resets (email, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE
		SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at`, 
		req.Email, token, expiresAt)
	
	if err != nil {
		user_logs.Create_user_log(tx, &adminID, "admin", "ADMIN_TRIGGER_RESET_FAILED", err.Error(), "FAILED", time.Since(start), &targetUserID)
		tx.Commit(ctx)
		http.Error(w, "Failed to generate reset token", http.StatusInternalServerError)
		return
	}

	// Audit Log
	user_logs.Create_user_log(tx, &adminID, "admin", "ADMIN_TRIGGER_RESET_SUCCESS", fmt.Sprintf("Triggered password reset for %s", req.Email), "SUCCESS", time.Since(start), &targetUserID)

	// Notifications
	notifications.Send_notification(adminID, tx, fmt.Sprintf("Triggered password reset for user: %s", req.Email), "Reset Triggered")
	notifications.Send_notification(targetUserID, tx, "An administrator has initiated a password reset for your account. Please check your email.", "Password Reset Triggered")

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send Email (Hardcoded recipient per project requirement)
	services.SendPasswordResetEmail("richardsambo94@gmail.com", token)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset triggered successfully"})
}
