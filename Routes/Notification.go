package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"ubaps/Db"
	middleware "ubaps/Middleware"
	"ubaps/utils"
)


func Notifications(w http.ResponseWriter, r *http.Request) {
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

	notifications, err := utils.ReciveNotifications(Db.DB, ctx, userID)
	if err != nil {
		log.Println("Error fetching notifications:", err)
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}

	if len(notifications) == 0 {
		json.NewEncoder(w).Encode([]any{})
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(notifications)
	if err != nil {
		http.Error(w, "Failed to encode notifications", http.StatusInternalServerError)
		return
	}
}

func NotificationCounter(w http.ResponseWriter, r *http.Request) {
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

	count, err := utils.NotificationCounter(Db.DB, ctx, userID)
	if err != nil {
		log.Println("Error fetching notificationCount:", err)
		http.Error(w, "Failed to fetch notificationCount", http.StatusInternalServerError)
		return
	}
    w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(count)
	if err != nil {
    http.Error(w,"Failed to send Counter",http.StatusInternalServerError)
	}
}
