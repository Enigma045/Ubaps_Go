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
	notifications "ubaps/Notifications"
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
		http.Error(w,"This is no a post methord",http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	tx,err := Db.DB.Begin(ctx)
	if err != nil{
		log.Println(err)
		http.Error(w,"Transction Failed to create",http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)

	

	err = json.NewDecoder(r.Body).Decode(&user)
    if err != nil{
		log.Println(err)
		http.Error(w,"Failed to get payload",http.StatusInternalServerError)
		return
	}
    
	userID,err := Handles.CreateUser(user.First,user.Last,user.Email,user.Phone,user.Password,user.Role,tx,true)
    if err != nil{
		log.Println(err)
		http.Error(w,"Failed to create user",http.StatusInternalServerError)
		return
	}
   log.Println(user)
	err = notifications.Send_notification(userID, tx, "Your account has been created.","Account Created")
	if err != nil {
		log.Println(err)
		http.Error(w,"Failed to create user notification",http.StatusInternalServerError)
		return
	}

	

	duration := time.Since(start)
	err = user_logs.Create_user_log(tx, &userID, user.Role, fmt.Sprintf("%s_ACCOUNT_CREATED", user.Role), fmt.Sprintf("user:%d", userID), "SUCCESS", duration)
	if err != nil {
		log.Println(err)
		http.Error(w,"Failed to create user log",http.StatusInternalServerError)
		return
	}

	
	err = tx.Commit(ctx)
	if err != nil{
		log.Println(err)
		http.Error(w,"Transction Failed to commit",http.StatusInternalServerError)
		return
	}
    w.WriteHeader(http.StatusOK)
	w.Write([]byte("User succefully created"))
}