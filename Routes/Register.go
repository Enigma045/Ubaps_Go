package Routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	user_logs "ubaps/Audit_logs"
	"ubaps/Db"
	"ubaps/Handles"
	notifications "ubaps/Notifications"
	"ubaps/services"
	"ubaps/utils"

	"github.com/jackc/pgx/v5"
)

type RegisterRequest struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	RegNumber string `json:"reg_number"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	start := time.Now()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Decode Error:", err)
		user_logs.Create_user_log(nil, nil, "student", "REGISTRATION_FAILED", "Invalid JSON payload", "FAILED", time.Since(start), nil)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	// --- Rolling 5-Year Eligibility Window Check ---
	parts := strings.Split(req.RegNumber, "-")
	if len(parts) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid registration number format"})
		return
	}
	regYear, err := strconv.Atoi(parts[3])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid year in registration number"})
		return
	}
	currentYear := time.Now().Year() % 100
	minYear := currentYear - 4
	if regYear < minYear || regYear > currentYear {
		fullMinYear := 2000 + minYear
		fullMaxYear := 2000 + currentYear
		errMsg := fmt.Sprintf("Registration is only open to students enrolled between %d and %d. Your enrolment year (%d) is outside this window.", fullMinYear, fullMaxYear, 2000+regYear)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
		return
	}

	// Start transaction
	ctx := context.Background()
	tx, err := Db.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Println("Begin Tx Error:", err)
		user_logs.Create_user_log(nil, nil, "student", "REGISTRATION_FAILED", "Database error", "FAILED", time.Since(start), nil)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	rollback := true
	defer func() {
		if rollback {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// Create User
	email := fmt.Sprintf("%s@unilia.ac.mw", req.RegNumber)
	userID, err := Handles.CreateUser(req.Name, req.Surname, email, req.Phone, req.Password, "student", tx, true)
	if err != nil {
		log.Println("CreateUser Error:", err)
		user_logs.Create_user_log(tx, nil, "student", "REGISTRATION_FAILED", err.Error(), "FAILED", time.Since(start), nil)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Notifications & Token
	notifications.Send_notification(userID, tx, "Your account has been created.", "Account Created")
	token, err := utils.GenerateVerificationToken(email, tx)
	if err != nil {
		log.Println("Token Generation Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate verification token"})
		return
	}

	// Success
	rollback = false
	services.SendVerificationEmail("richardsambo94@gmail.com", token) // Non-fatal
	
	duration := time.Since(start)
	user_logs.Create_user_log(tx, &userID, "student", "STUDENT_ACCOUNT_CREATED", fmt.Sprintf("user:%d", userID), "SUCCESS", duration, &userID)
	
	if adminIDs, adminErr := Handles.GetAdmins(tx); adminErr == nil {
		notifications.BroadcastNotification(adminIDs, tx, fmt.Sprintf("A new student user (%s) has registered.", email), "New User Registration")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Registration successful. Check your email.",
		"email":   email,
	})
}
