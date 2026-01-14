package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func CreateSessionTx(
	ctx context.Context,
	w http.ResponseWriter,
	tx pgx.Tx,
	userID int,
) error {
	sessionID := uuid.New()
	expires := time.Now().Add(2 * time.Minute)

	_, err := tx.Exec(ctx, `
	INSERT INTO sessions (session_id,user_id, expires_at)
	VALUES ($1,$2,$3) `,
		sessionID, userID, expires)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID.String(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expires,
		Path:     "/"})
	return nil
}
