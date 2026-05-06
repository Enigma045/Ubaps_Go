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
	notifications "ubaps/Notifications"
	"ubaps/services"
)
type User struct{

	First string `json:"first_name"`
	Last string `json:"last_name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Password string `json:"password"`
	Role string `json:"role"`
}

func CreateUser(w http.ResponseWriter,r *http.Request) {
	start := time.Now()
	var user User
    if r.Method != http.MethodPost{
		http.Error(w,"This is not a post method",http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	tx,err := Db.DB.Begin(ctx)
	if err != nil{
		log.Println(err)
		http.Error(w,"Transaction Failed to create",http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)

	

	performerID, _ := middleware.UserIDFromContext(ctx)
	var performerIDPtr *int64
	if performerID != 0 {
		performerIDPtr = &performerID
	}

	err = json.NewDecoder(r.Body).Decode(&user)
    if err != nil{
		log.Println(err)
		user_logs.Create_user_log(tx, performerIDPtr, "admin", "USER_CREATION_FAILED", "Invalid JSON payload", "FAILED", time.Since(start), nil)
		http.Error(w,"Failed to get payload",http.StatusInternalServerError)
		return
	}
    
	userID,err := Handles.CreateUser(user.First,user.Last,user.Email,user.Phone,user.Password,user.Role,tx,true)
    if err != nil{
		log.Println(err)
		user_logs.Create_user_log(tx, performerIDPtr, "admin", "USER_CREATION_FAILED", fmt.Sprintf("Email:%s Error:%s", user.Email, err.Error()), "FAILED", time.Since(start), nil)
		http.Error(w,"Failed to create user",http.StatusInternalServerError)
		return
	}
   log.Println(user)
	err = notifications.Send_notification(userID, tx, "Your account has been created.","Account Created")
	if err != nil {
		log.Println(err)
		user_logs.Create_user_log(tx, performerIDPtr, user.Role, "USER_NOTIFICATION_FAILED", "Notification delivery failed during creation", "FAILED", time.Since(start), &userID)
		http.Error(w,"Failed to create user notification",http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	err = user_logs.Create_user_log(tx, performerIDPtr, "admin", "USER_ACCOUNT_CREATED", fmt.Sprintf("user:%d, role:%s", userID, user.Role), "SUCCESS", duration, &userID)
	if err != nil {
		log.Println(err)
	}

	// Notify all admins
	if adminIDs, err := Handles.GetAdmins(tx); err == nil {
		notifications.BroadcastNotification(adminIDs, tx, fmt.Sprintf("A new user (%s) has been successfully created with role: %s", user.Email, user.Role), "User Created")
	}
	
	err = tx.Commit(ctx)
	if err != nil {
		log.Println(err)
		user_logs.Create_user_log(nil, performerIDPtr, "admin", "USER_CREATION_COMMIT_FAILED", user.Email, "FAILED", time.Since(start), &userID)
		http.Error(w, "Transaction Failed to commit", http.StatusInternalServerError)
		return
	}

	// Send Welcome Email
	err = services.SendWelcomeEmail(user.Email, user.Password, user.Role)
	if err != nil {
		log.Println("Failed to send welcome email:", err)
		user_logs.Create_user_log(nil, performerIDPtr, "admin", "WELCOME_EMAIL_FAILED", fmt.Sprintf("User: %s, Error: %s", user.Email, err.Error()), "FAILED", time.Since(start), &userID)
	} else {
		user_logs.Create_user_log(nil, performerIDPtr, "admin", "WELCOME_EMAIL_SENT", fmt.Sprintf("Welcome email sent to %s", user.Email), "SUCCESS", time.Since(start), &userID)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User successfully created and welcome email sent"})
}