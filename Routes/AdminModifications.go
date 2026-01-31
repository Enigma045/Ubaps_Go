package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
	middleware "ubaps/Middleware"
	"ubaps/utils"
)
type emailRequest struct {
	Email string `json:"email"`
}

func Getuserdetails(w http.ResponseWriter,r *http.Request){
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

	count, err := utils.ReciveDetails(Db.DB, ctx, userID)
	if err != nil {
		log.Println("Error fetching userDetails:", err)
		http.Error(w, "Failed to fetch userDetails", http.StatusInternalServerError)
		return
	}
	log.Println("Notification count:", count)
    w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(count)
	if err != nil {
    http.Error(w,"Failed to send Counter",http.StatusInternalServerError)
	}
}

func DeleteAccount(w http.ResponseWriter,r *http.Request){
    var emailreq emailRequest

	if r.Method != http.MethodPost{
    http.Error(w,"wrong Method",http.StatusMethodNotAllowed)
	}

	ctx := r.Context()

	tx,err := Db.DB.Begin(ctx)
    if err != nil {
		log.Println(err)
		http.Error(w,"Failed to create Transction",http.StatusInternalServerError)
	    return
	}
	defer tx.Rollback(ctx)

	err = json.NewDecoder(r.Body).Decode(&emailreq)
	if err != nil {
		log.Println(err)
		http.Error(w,"Invalid JSON",http.StatusBadRequest)
		return
	}
	
	userid,err := Handles.GetUserIDByEmail(emailreq.Email,tx)
	if err != nil {
		log.Println(err)
		http.Error(w,"Failed to get userid",http.StatusInternalServerError)
	    return
	}

	err = utils.DeleteUser(tx,ctx,userid)
	if err != nil {
		log.Println(err)
		http.Error(w,"Invalid JSON",http.StatusBadRequest)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You have succefully deleted the user"))

}
