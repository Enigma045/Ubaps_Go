package Routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"ubaps/Db"
)

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Invalid token", 400)
		return
	}

	var email string
	var expires time.Time

	err := Db.DB.QueryRow(
		context.Background(),
		`SELECT email, expires_at FROM email_verifications WHERE token=$1`,
		token,
	).Scan(&email, &expires)

	if err != nil || time.Now().After(expires) {
		http.Error(w, "Token expired or invalid", 400)
		return
	}

	Db.DB.Exec(
		context.Background(),
		`UPDATE users SET is_verified=true WHERE email=$1`,
		email,
	)

	Db.DB.Exec(
		context.Background(),
		`DELETE FROM email_verifications WHERE token=$1`,
		token,
	)

	fmt.Println("success5")

	w.Write([]byte("Email verified successfully"))
}
