package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"ubaps/Db"
	"github.com/google/uuid"
)

/*
|--------------------------------------------------------------------------
| Context Keys (type-safe)
|--------------------------------------------------------------------------
*/
type contextKey string

const (
	userIDKey contextKey = "user_id"
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
			log.Println("AUTH_DEBUG: No session_id cookie found in request:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("AUTH_DEBUG: Received session_id cookie:", cookie.Value)

		// Parse sessionID as UUID
		parsedID, err := uuid.Parse(cookie.Value)
		if err != nil {
			log.Println("AUTH_DEBUG: Invalid UUID in cookie:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var (
			userID    int64
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
			`,
			parsedID,
		).Scan(&userID, &role, &expiresAt)
		log.Println("AUTH_DEBUG: Query complete.")

		if err != nil {
			log.Println("AUTH_DEBUG: Session lookup failed for ID", parsedID, ":", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		currentTime := time.Now().UTC()
		log.Println("AUTH_DEBUG: Session found. Expiration (DB):", expiresAt.Format(time.RFC3339), "Current Time (UTC):", currentTime.Format(time.RFC3339))

		if expiresAt.Before(currentTime) {
			log.Println("AUTH_DEBUG: Session is EXPIRED according to Go UTC comparison.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Println("AUTH_DEBUG: Session is VALID. UserID:", userID, "Role:", role)

		// 🔄 Sliding session refresh
		if time.Until(expiresAt) < refreshThreshold {
			newExpiry := time.Now().Add(sessionDuration)

			_, err := Db.DB.Exec(
				r.Context(),
				`
				UPDATE sessions
				SET expires_at = $1
				WHERE session_id = $2
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
					Secure:   false,
					Expires:  newExpiry,
					SameSite: http.SameSiteLaxMode,
				})
			}
		}

		// ✅ Store values in request context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, roleKey, role)
		if id, ok := ctx.Value(userIDKey).(int64); !ok {
			log.Println("FAILED to store userID in context")
		} else {
			log.Println("Stored userID in context:", id)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*
|--------------------------------------------------------------------------
| RequireRole Middleware
|--------------------------------------------------------------------------
*/
func RequireRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value(roleKey).(string)
			if !ok {
				http.Error(w, "Forbidden: role not found", http.StatusForbidden)
				return
			}

			// Check if user's role matches any allowed role
			for _, allowedRole := range requiredRoles {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: not your role", http.StatusForbidden)
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
			`DELETE FROM sessions WHERE session_id = $1`,
			cookie.Value,
		)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/Login", http.StatusSeeOther)
}

/*
|--------------------------------------------------------------------------
| Helper Functions (optional but clean)
|--------------------------------------------------------------------------
*/
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey).(string)
	return role, ok
}
