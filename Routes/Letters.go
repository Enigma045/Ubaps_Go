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
	"ubaps/services"
	"fmt"
)

type LetterMetadata struct {
	ID                 int       `json:"id"`
	StudentName        string    `json:"studentName"`
	RegistrationNumber string    `json:"registrationNumber"`
	LetterName         string    `json:"letterName"`
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

	letterName := header.Filename

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
		user_logs.Create_user_log(tx, &userId, "student", "STUDENT_LETTER_ALREADY_SENT", fmt.Sprintf("user:%d", userId), "FAILED", time.Since(start), &userId)
		
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
		INSERT INTO letters (user_id, letter, letter_name, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, query, userId, content, letterName, time.Now())
	if err != nil {
		log.Println("Database error in SubmitLetter:", err)
		http.Error(w, "Failed to save letter", http.StatusInternalServerError)
		return
	}

	// Success path
	notifications.Send_notification(userId, tx, "You have successfully sent the student letter.", "Letter Submitted")
	user_logs.Create_user_log(tx, &userId, "student", "STUDENT_LETTER_SUBMITTED", fmt.Sprintf("user:%d", userId), "SUCCESS", time.Since(start), &userId)

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
			l.letter_name,
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

		if err := rows.Scan(&l.ID, &l.StudentName, &regNum, &l.LetterName, &createdAt); err != nil {
			log.Println("Scan error in GetLettersList:", err)
			continue
		}

		if regNum != nil {
			l.RegistrationNumber = *regNum
		} else {
			l.RegistrationNumber = "N/A"
		}
		
		l.DateSubmitted = createdAt.Format("2006-01-02 15:04")
		if l.LetterName != "" {
			l.LetterType = l.LetterName
		} else {
			l.LetterType = "Student Letter"
		}
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
	var letterName string
	err := Db.DB.QueryRow(r.Context(), "SELECT letter, letter_name FROM letters WHERE letters_id = $1", id).Scan(&content, &letterName)
	if err != nil {
		http.Error(w, "Letter not found", http.StatusNotFound)
		return
	}

	// Fallback if letterName is empty
	if letterName == "" {
		letterName = "letter_" + id
	}

	// Detection of content type (simplified)
	contentType := http.DetectContentType(content)
	
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", letterName))
	w.Write(content)
}

func SendLetter(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing letter ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var content []byte
	var letterName string
	var studentUserId int64

	// 1. Get letter data
	err = tx.QueryRow(ctx, "SELECT letter, letter_name, user_id FROM letters WHERE letters_id = $1", id).Scan(&content, &letterName, &studentUserId)
	if err != nil {
		log.Println("Error fetching letter for sending:", err)
		http.Error(w, "Letter not found", http.StatusNotFound)
		return
	}

	// 2. Get benefactor email for this student
	var benefactorEmail string
	var benefactorName string
	query := `
		SELECT bs.benefactor_email, bs.benefactor_name
		FROM bursary_schemes bs
		JOIN applications a ON a.scheme_id = bs.scheme_id
		WHERE a.user_id = $1
	`
	err = tx.QueryRow(ctx, query, studentUserId).Scan(&benefactorEmail, &benefactorName)
	if err != nil {
		log.Println("Error fetching benefactor email for student:", studentUserId, err)
		http.Error(w, "Benefactor information not found for this student. Ensure the student is assigned to a scheme.", http.StatusBadRequest)
		return
	}

	// 3. Send email
	subject := "New Student Letter Received: " + letterName
	body := fmt.Sprintf(`
		<h3>Dear %s,</h3>
		<p>You have received a new letter from a student under your bursary scheme.</p>
		<p>Please find the attached document: <strong>%s</strong></p>
		<br>
		<p>Best regards,<br>UBAPS System</p>
	`, benefactorName, letterName)

	err = services.SendEmailWithAttachment(benefactorEmail, subject, body, letterName, content)
	if err != nil {
		log.Println("Failed to send email to benefactor:", err)
		user_logs.Create_user_log(tx, nil, "staff", "FORWARD_LETTER_FAILED", fmt.Sprintf("LetterID:%s, Benefactor:%s", id, benefactorEmail), "FAILED", time.Since(start), &studentUserId)
		http.Error(w, "Failed to send email to benefactor", http.StatusInternalServerError)
		return
	}

	// Notify student that their letter was forwarded
	notifications.Send_notification(studentUserId, tx, fmt.Sprintf("Your student letter (%s) has been successfully forwarded to your benefactor (%s).", letterName, benefactorName), "Letter Forwarded")

	// Audit Log
	staffID, _ := middleware.UserIDFromContext(ctx)
	user_logs.Create_user_log(tx, &staffID, "staff", "STUDENT_LETTER_FORWARDED", fmt.Sprintf("LetterID:%s, Student:%d, Benefactor:%s", id, studentUserId, benefactorEmail), "SUCCESS", time.Since(start), &studentUserId)

	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Letter sent to benefactor successfully"))
}
