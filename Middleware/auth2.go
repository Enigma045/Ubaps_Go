package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"ubaps/Db"
)

/*
|--------------------------------------------------------------------------
| Context Keys (type-safe)
|--------------------------------------------------------------------------
*/
type contextKey string

const (
	userIDKey contextKey = "userID"
	roleKey   contextKey = "role"
)

/*
|--------------------------------------------------------------------------
| Session Config
|--------------------------------------------------------------------------
*/
const (
	sessionDuration  = 30 * time.Minute
	refreshThreshold = 10 * time.Minute
)

/*
|--------------------------------------------------------------------------
| RequireAuth Middleware
|--------------------------------------------------------------------------
| - Validates session
| - Loads user ID + role
| - Refreshes session if close to expiry
| - Injects values into request context
|--------------------------------------------------------------------------
*/
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Session not found Unauthorized", http.StatusUnauthorized)
			return
		}

		var (
			userID    int
			role      string
			expiresAt time.Time
		)

		err = Db.DB.QueryRow(
			r.Context(),
			`
			SELECT s.user_id, u.user_type, s.expires_at
			FROM sessions s
			JOIN users u ON u.user_id = s.user_id
			WHERE s.session_id = $1
			  AND s.expires_at > NOW()
			`,
			cookie.Value,
		).Scan(&userID, &role, &expiresAt)

		if err != nil {
			log.Println(cookie.Value)
			http.Error(w, "Session not found in database Unauthorized", http.StatusUnauthorized)
			return
		}

		// 🔄 Sliding session refresh
		if time.Until(expiresAt) < refreshThreshold {
			newExpiry := time.Now().Add(sessionDuration)

			_, err := Db.DB.Exec(
				r.Context(),
				`
				UPDATE sessions
				SET expires_at = $1
				WHERE id = $2
				`,
				newExpiry,
				cookie.Value,
			)

			if err == nil {
				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    cookie.Value,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					Expires:  newExpiry,
					SameSite: http.SameSiteStrictMode,
				})
			}
		}

		// Store values in context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, roleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*
|--------------------------------------------------------------------------
| RequireRole Middleware
|--------------------------------------------------------------------------
*/
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value(roleKey).(string)
			if !ok || role != requiredRole {
				http.Error(w, "Forbidden not you Role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

/*
|--------------------------------------------------------------------------
| RequireAnyRole Middleware (recommended)
|--------------------------------------------------------------------------
*/
func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{})
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value(roleKey).(string)
			if !ok {
				http.Error(w, "Forbidden Failed to retrieve Role in database", http.StatusForbidden)
				return
			}

			if _, ok := allowed[role]; !ok {
				http.Error(w, "Forbidden Failed to match Role with Database", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

/*
|--------------------------------------------------------------------------
| Logout Handler
|--------------------------------------------------------------------------
*/
func Logout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_id")
	if err == nil {
		_, _ = Db.DB.Exec(
			r.Context(),
			`DELETE FROM sessions WHERE id = $1`,
			cookie.Value,
		)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out"))
}

/*
|--------------------------------------------------------------------------
| Helper Functions (optional but clean)
|--------------------------------------------------------------------------
*/
func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey).(string)
	return role, ok
}
