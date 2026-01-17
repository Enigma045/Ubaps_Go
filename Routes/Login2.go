package Routes

import (
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

	var (
		userID      int64
		hash        string
		is_verified bool
		role        string
	)
	err = tx.QueryRow(ctx, `
	SELECT user_id, password_hash, is_verified , user_type
	FROM users WHERE email = $1
	`,
		r.FormValue("email")).Scan(&userID, &hash, &is_verified, &role)
	if err != nil {

		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if !is_verified {
		token, err := utils.GenerateVerificationToken(r.FormValue("email"), tx)
		if err != nil {
			http.Error(w, "Token Generation Failed", http.StatusInternalServerError)

		}

		err = services.SendVerificationEmail(r.FormValue("email"), token)
		if err != nil {
			http.Error(w, "Failed to send verification Email", http.StatusInternalServerError)
		}

		// ✅ COMMIT BEFORE RETURN
		if err := tx.Commit(ctx); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		http.Error(w, "This Account is not Verified please verify using your school account", http.StatusForbidden)
		return
	}
	if bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(r.FormValue("password")),
	) != nil {
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}

	err = utils.CreateSessionTx(ctx, w, tx, int(userID))
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	err = utils.FirstFill(ctx, role, userID, tx)
	if err != nil {
		log.Println("First Insertion Failed")
		http.Error(w, "Server error", 500)
		return
	}
	duration := time.Since(start)
	user_logs.Create_user_log(tx, &userID, role, "LOGGED_IN_ACCOUNT", fmt.Sprintf("user:%d", userID), "SUCCESS", duration)

	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Login successful"))
}
