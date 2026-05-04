package Routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"ubaps/Db"
	middleware "ubaps/Middleware"
	notifications "ubaps/Notifications"
	user_logs "ubaps/Audit_logs"
	"ubaps/Handles"
	"fmt"
)

type LetterMetadata struct {
	ID                 int       `json:"id"`
	StudentName        string    `json:"studentName"`
	RegistrationNumber string    `json:"registrationNumber"`
	LetterType         string    `json:"letterType"`
	DateSubmitted      string    `json:"dateSubmitted"`
	Status             string    `json:"status"`
}

func SubmitLetter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "File size too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("letter")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type by extension
	ext := strings.ToLower(header.Filename[strings.LastIndex(header.Filename, ".")+1:])
	if ext != "pdf" && ext != "doc" && ext != "docx" {
		http.Error(w, "Only PDF and Word documents are allowed", http.StatusBadRequest)
		return
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	start := time.Now()
	ctx := r.Context()
	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	// Check if letter already exists
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM letters WHERE user_id = $1)", userId).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists {
		// Notify student they already sent the letter
		notifications.Send_notification(userId, tx, "You have already sent the student letter.", "Letter Already Sent")
		// Log failure
		user_logs.Create_user_log(tx, &userId, "student", "STUDENT_LETTER_ALREADY_SENT", fmt.Sprintf("user:%d", userId), "FAILED", time.Since(start))
		
		if err := tx.Commit(ctx); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You have already sent the student letter"))
		return
	}

	// Insert into letters table
	query := `
		INSERT INTO letters (user_id, letter, created_at)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, query, userId, content, time.Now())
	if err != nil {
		log.Println("Database error in SubmitLetter:", err)
		http.Error(w, "Failed to save letter", http.StatusInternalServerError)
		return
	}

	// Success path
	notifications.Send_notification(userId, tx, "You have successfully sent the student letter.", "Letter Submitted")
	user_logs.Create_user_log(tx, &userId, "student", "STUDENT_LETTER_SUBMITTED", fmt.Sprintf("user:%d", userId), "SUCCESS", time.Since(start))

	// Get user email for broadcast
	email, _ := Handles.GetEmailByUserID(userId, tx)

	// Broadcast to users of different types (excluding students and admin)
	if staffIDs, err := Handles.GetUserIDsOfDifferentTypes(tx, "student"); err == nil {
		notifications.BroadcastNotification(staffIDs, tx, fmt.Sprintf("A student (%s) has submitted a letter.", email), "New Letter Received")
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Letter submitted successfully"))
}

func GetLettersList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Join with users and applications to get metadata
	query := `
		SELECT 
			l.letters_id,
			u.name || ' ' || u.surname as full_name,
			a.registration_number,
			l.created_at
		FROM letters l
		JOIN users u ON u.user_id = l.user_id
		LEFT JOIN applications a ON a.user_id = l.user_id
		ORDER BY l.created_at DESC
	`

	rows, err := Db.DB.Query(ctx, query)
	if err != nil {
		log.Println("Database error in GetLettersList:", err)
		http.Error(w, "Failed to fetch letters", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var letters []LetterMetadata
	for rows.Next() {
		var l LetterMetadata
		var createdAt time.Time
		var regNum *string

		if err := rows.Scan(&l.ID, &l.StudentName, &regNum, &createdAt); err != nil {
			log.Println("Scan error in GetLettersList:", err)
			continue
		}

		if regNum != nil {
			l.RegistrationNumber = *regNum
		} else {
			l.RegistrationNumber = "N/A"
		}
		
		l.DateSubmitted = createdAt.Format("2006-01-02 15:04")
		l.LetterType = "Thank You Letter"
		l.Status = "received"
		letters = append(letters, l)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(letters)
}

func DownloadLetter(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing letter ID", http.StatusBadRequest)
		return
	}

	var content []byte
	err := Db.DB.QueryRow(r.Context(), "SELECT letter FROM letters WHERE letters_id = $1", id).Scan(&content)
	if err != nil {
		http.Error(w, "Letter not found", http.StatusNotFound)
		return
	}

	// Detection of content type (simplified)
	contentType := http.DetectContentType(content)
	
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=letter_"+id)
	w.Write(content)
}
