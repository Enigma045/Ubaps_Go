package Routes

import (
	"log"
	"net/http"
	"ubaps/Db"
	middleware "ubaps/Middleware"
	"ubaps/utils"
)

func Scheme_Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx := r.Context()
	// Role, rcheck := middleware.RoleFromContext(ctx)
	// if rcheck != true {
	// 	log.Println("Failed to take UserId")
	// }
	// Begin transaction
	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx) // Rollback will be ignored if commit succeeds
	//
	userId, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	log.Println("User ID:", userId)
	// Parse form safely
	formData, err := utils.Formdata(r)
	if err != nil {
		log.Println("Formdata Error:", err)
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	log.Println("FormData received:", formData)

	// Execute DB operation safely
	if err := utils.Scheme_Operations(tx, ctx, formData, userId); err != nil {
		log.Println("DB Operation Failed:", err)
		http.Error(w, "Database operation failed", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You have successfully submitted the application form"))
}
