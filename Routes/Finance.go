package Routes

import (
	"fmt"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
	"ubaps/utils"
)

func Fees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx := r.Context()
	tx, err := Db.DB.Begin(ctx)

	if err != nil {
		http.Error(w, "Database Failed transctions Error ", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)
	reg := r.FormValue("student_id")
	email := fmt.Sprintf("%s@unilia.ac.mw", reg)

	data, err := utils.Formdata(r)
	if err != nil {
		http.Error(w, "Formdate Error ", http.StatusInternalServerError)
		return
	}
	log.Println(email)
	student_id, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get student_id ", http.StatusInternalServerError)
		return
	}
	err = utils.Finance_Operations(tx, ctx, data, student_id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get student_id ", http.StatusInternalServerError)
		return
	}
	tx.Commit(ctx)
	if err != nil {
		http.Error(w, "Database Failed to commit Error ", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Fees statement successufuly sent"))
}

func Request_Statement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx := r.Context()

	// Parse student_id (registration number) from request body (JSON)
	var payload struct {
		Reg string `json:"reg"`
	}
	if err := utils.DecodeJSON(r, &payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.Reg == "" {
		http.Error(w, "Student registration number is required", http.StatusBadRequest)
		return
	}

	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	email := fmt.Sprintf("%s@unilia.ac.mw", payload.Reg)
	studentID, err := Handles.GetUserIDByEmail(email, tx)
	if err != nil {
		log.Println("Error resolving student ID:", err)
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	if err := utils.Request_Finance_Statement(tx, ctx, studentID); err != nil {
		log.Println("Error requesting statement:", err)
		http.Error(w, "Failed to request statement", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Statement request sent successfully"}`))
}
