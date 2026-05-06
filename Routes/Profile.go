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
	"ubaps/Notifications"
)

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role, _ := middleware.RoleFromContext(ctx)

	var req struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
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

	err = Handles.UpdateUserProfile(ctx, tx, userID, req.Name, req.Surname, req.Email, req.Phone)
	if err != nil {
		log.Println("Error updating profile:", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Audit Log
	logMsg := fmt.Sprintf("Updated personal info: Name=%s, Surname=%s, Email=%s, Phone=%s", req.Name, req.Surname, req.Email, req.Phone)
	user_logs.Create_user_log(tx, &userID, role, "PROFILE_UPDATED", logMsg, "SUCCESS", time.Since(start), &userID)

	// In-app Notification
	notifications.Send_notification(userID, tx, "You have successfully updated your profile details.", "Profile Updated")

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role, _ := middleware.RoleFromContext(ctx)

	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
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

	// 1. Verify old password
	var hash string
	err = tx.QueryRow(ctx, "SELECT password_hash FROM users WHERE user_id = $1", userID).Scan(&hash)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !Handles.CheckPassword(hash, req.OldPassword) {
		user_logs.Create_user_log(tx, &userID, role, "PASSWORD_CHANGE_FAILED", "Incorrect old password", "FAILED", time.Since(start), &userID)
		tx.Commit(ctx) // Commit the failed log
		http.Error(w, "Incorrect old password", http.StatusUnauthorized)
		return
	}

	// 2. Update to new password
	err = Handles.UpdateUserPassword(ctx, tx, userID, req.NewPassword)
	if err != nil {
		log.Println("Error updating password:", err)
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Audit Log
	user_logs.Create_user_log(tx, &userID, role, "PASSWORD_CHANGED", "User successfully changed their password", "SUCCESS", time.Since(start), &userID)

	// In-app Notification
	notifications.Send_notification(userID, tx, "Your password has been changed successfully.", "Password Changed")

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}
