package Routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	user_logs "ubaps/Audit_logs"
	"ubaps/Db"
	"ubaps/services"
	"ubaps/utils"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	start := time.Now()
	ctx := r.Context()
	log.Println(r.FormValue("email"))
	log.Println(r.FormValue("password"))
	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	defer tx.Rollback(ctx)

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload"})
			return
		}
	} else {
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
	}

	var (
		userID      int64
		hash        string
		is_verified bool
		role        string
		is_active   bool
	)

	err = tx.QueryRow(ctx, `
	SELECT user_id, password_hash, is_verified, user_type, is_active
	FROM users WHERE email = $1
	`, req.Email).Scan(&userID, &hash, &is_verified, &role, &is_active)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid email or password"})
		return
	}

	if !is_active {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "Your account is deactivated. Please contact the administrator."})
		return
	}

	if !is_verified {
		token, err := utils.GenerateVerificationToken(req.Email, tx)
		if err == nil {
			services.SendVerificationEmail(req.Email, token)
		}
		
		if err := tx.Commit(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Server error"})
			return
		}

		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "This account is not verified. Please verify using your school account."})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid password"})
		return
	}

	err = utils.CreateSessionTx(ctx, w, tx, int(userID))
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	err = utils.FirstFill(ctx, role, userID, tx)
	if err != nil {
		log.Println("First Insertion Failed",err)
		http.Error(w, "Server error", 500)
		return
	}
	duration := time.Since(start)
	user_logs.Create_user_log(tx, &userID, role, "LOGGED_IN_ACCOUNT", fmt.Sprintf("user:%d", userID), "SUCCESS", duration, &userID)

	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"role":role,
	})
	if err != nil {
		log.Println("Failed to send role to front end:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}
