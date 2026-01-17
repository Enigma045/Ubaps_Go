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
		log.Println("Error parsing JSON:", err)
	}

	// Extract a value
	reg, ok := jsonData[filter].(string)
	if !ok {
		log.Println(fmt.Sprintf("%s not found or not a string", filter))
	}

	return reg, err
}

func Filter(body []byte, info [5]string, tx pgx.Tx) string {
	//loop & array hell

	var details [len(info)]string

	for i := 0; i < len(info); i++ {
		value, err := Contains(body, info[i])
		if err != nil {
			log.Println("failed to extract string:", err)
			continue
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
	Handles.CreateUser(name, surname, email, phone, password, "student", tx)
	return email
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read request body"})
		return
	}
	defer r.Body.Close()

	// Start transaction
	tx, err := Db.DB.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
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
	reqEmail := Filter(body, info, tx) // Filter should call Handles.CreateUser inside tx and return email
	if reqEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse user info"})
		return
	}

	// Insert notification
	userID, err := Handles.GetUserIDByEmail(reqEmail, tx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve new user ID"})
		return
	}
	err = notifications.User_Created(userID, tx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create notification"})
		return
	}

	// Generate verification token
	token, err := utils.GenerateVerificationToken(reqEmail, tx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		return
	}

	// Send verification email (outside transaction)
	err = services.SendVerificationEmail(reqEmail, token)
	if err != nil {
		log.Println("Email send error:", err)
	}
	fmt.Println("success7")
	// Respond with valid JSON
	//User_logs
	duration := time.Since(start)
	user_logs.Create_user_log(tx, &userID, "student", "STUDENT_ACCOUNT_CREATED", fmt.Sprintf("user:%d", userID), "SUCCESS", duration)
	//i wonder whatt happens if the acountis not verified
	//User_logs
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Registration successful. Check your email.",
		"email":   reqEmail,
	})
}
