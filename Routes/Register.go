package Routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	user_logs "ubaps/Audit_logs"
	"ubaps/Db"
	"ubaps/Handles"
	notifications "ubaps/Notifications"
	"ubaps/services"
	"ubaps/utils"

	"github.com/jackc/pgx/v5"
)

func Contains(body []byte, filter string) (string, error) {
	// Parse JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	// Extract a value
	val, ok := jsonData[filter]
	if !ok {
		return "", fmt.Errorf("%s not found", filter)
	}

	reg, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("%s is not a string", filter)
	}

	return reg, nil
}

func Filter(body []byte, info [5]string, tx pgx.Tx) (int64, string, error) {
	//loop & array hell
	var details [5]string

	for i := 0; i < len(info); i++ {
		value, err := Contains(body, info[i])
		if err != nil {
			return 0, "", fmt.Errorf("failed to extract %s: %w", info[i], err)
		}
		details[i] = value
	}

	// Assign extracted values
	name := details[0]
	surname := details[1]
	phone := details[2]
	password := details[3]
	reg := details[4]

	// Build email
	email := fmt.Sprintf("%s@unilia.ac.mw", reg)
	log.Println("name:", name)
	log.Println("surname:", surname)
	log.Println("phone:", phone)
	log.Println("password:", password)
	log.Println("email:", email)
	//fall loop and array
	fmt.Println("success1")
	userID, err := Handles.CreateUser(name, surname, email, phone, password, "student", tx, true)
	if err != nil {
		log.Println("CreateUser Error:", err)
		return 0, "", err
	}
	return userID, email, nil
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	start := time.Now()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Read Body Error:", err)
		user_logs.Create_user_log(nil, nil, "student", "REGISTRATION_FAILED", "Failed to read request body", "FAILED", time.Since(start), nil)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read request body"})
		return
	}
	defer r.Body.Close()

	// Start transaction
	tx, err := Db.DB.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		log.Println("Begin Tx Error:", err)
		user_logs.Create_user_log(nil, nil, "student", "REGISTRATION_FAILED", "Failed to start database transaction", "FAILED", time.Since(start), nil)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to start DB transaction"})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	// Extract user info from JSON
	info := [5]string{"name", "surname", "phone", "password", "reg_number"}
	userID, reqEmail, err := Filter(body, info, tx)
	if err != nil {
		user_logs.Create_user_log(tx, nil, "student", "REGISTRATION_FAILED", err.Error(), "FAILED", time.Since(start), nil)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Insert notification
	err = notifications.Send_notification(userID, tx, "Your account has been created.", "Account Created")
	if err != nil {
		log.Println("Send_notification Error:", err)
		user_logs.Create_user_log(tx, &userID, "student", "REGISTRATION_NOTIFICATION_FAILED", "Failed to send welcome notification", "FAILED", time.Since(start), &userID)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create notification"})
		return
	}

	// Generate verification token
	token, err := utils.GenerateVerificationToken(reqEmail, tx)
	if err != nil {
		log.Println("Token Generation Error:", err)
		user_logs.Create_user_log(tx, &userID, "student", "REGISTRATION_FAILED", "Token generation failure", "FAILED", time.Since(start), &userID)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		return
	}

	// Send verification email (outside transaction logic, non-fatal)
	if mailErr := services.SendVerificationEmail("richardsambo94@gmail.com", token); mailErr != nil {
		log.Println("Email send error:", mailErr)
		user_logs.Create_user_log(tx, &userID, "student", "REGISTRATION_EMAIL_FAILED", "Verification email delivery failed", "FAILED", time.Since(start), &userID)
	}

	fmt.Println("success7")
	// Respond with valid JSON
	//User_logs
	duration := time.Since(start)
	user_logs.Create_user_log(tx, &userID, "student", "STUDENT_ACCOUNT_CREATED", fmt.Sprintf("user:%d", userID), "SUCCESS", duration, &userID)
	//i wonder whatt happens if the acountis not verified

	// Notify all admins (non-fatal)
	if adminIDs, adminErr := Handles.GetAdmins(tx); adminErr == nil {
		notifications.BroadcastNotification(adminIDs, tx, fmt.Sprintf("A new student user (%s) has registered.", reqEmail), "New User Registration")
	} else {
		log.Println("GetAdmins Error:", adminErr)
	}

	//User_logs
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Registration successful. Check your email.",
		"email":   reqEmail,
	})
}
