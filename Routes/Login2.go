package Routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	user_logs "ubaps/Audit_logs"
	"ubaps/Db"
	"ubaps/services"
	"ubaps/utils"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := context.Background()
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
	)
	err = tx.QueryRow(ctx, `
	SELECT user_id, password_hash, is_verified 
	FROM users WHERE email = $1
	`,
		r.FormValue("email")).Scan(&userID, &hash, &is_verified)
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
	duration := time.Since(start)
	user_logs.Create_user_log(tx, &userID, "student", "LOGGED_IN_ACCOUNT", fmt.Sprintf("user:%d", userID), "SUCCESS", duration)
	tx.Commit(ctx)
	w.Write([]byte("Login successful"))
}
func Login_page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/public/login.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}
