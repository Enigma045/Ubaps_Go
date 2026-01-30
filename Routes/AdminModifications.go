package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"ubaps/Db"
	middleware "ubaps/Middleware"
	"ubaps/utils"
)

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
