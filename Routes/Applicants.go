package Routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
	middleware "ubaps/Middleware"
	"ubaps/utils"
)

func Applicants(w http.ResponseWriter, r *http.Request) {
	var Pplicants []string
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	
	err := json.NewDecoder(r.Body).Decode(&Pplicants)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, limit, offset := utils.GetPaginationParams(r)

	results, total, err := Handles.Applicants(Db.DB, ctx, Pplicants, limit, offset)
	if err != nil {
		log.Println("Error retriving applicants from database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(utils.PaginatedResponse{
		Data:  results,
		Total: total,
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		log.Println("Error encording applicants to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type activeStudent struct {
	StudentName string `json:"name"`
	StudentID   string `json:"id"`
}

func ConsiderStudent(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content", "application/json")
	ctx := r.Context()

   tx, err := Db.DB.Begin(ctx)
   if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)
	

	var student activeStudent
	err = json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
    email := fmt.Sprintf("%s@unilia.ac.mw",student.StudentID)

	UserId, err := Handles.GetUserIDByEmail(email,tx)
	if err != nil {
		log.Println("Error getting user ID from email:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Role,okay := middleware.RoleFromContext(ctx)
	if okay != true{
		log.Println("Error Retrieving Role from Context", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	results, err := Handles.ConsiderStudent(tx, ctx, UserId, Role)
	if err != nil {
		log.Println("Error Considering applicants from database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit(ctx)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Println("Error encording applicants to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func RejectStudent(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content", "application/json")
	ctx := r.Context()

   tx, err := Db.DB.Begin(ctx)
   if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer tx.Rollback(ctx)

	var student activeStudent
	err = json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
    email := fmt.Sprintf("%s@unilia.ac.mw",student.StudentID)

	UserId, err := Handles.GetUserIDByEmail(email,tx)
	if err != nil {
		log.Println("Error getting user ID from email:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Role,okay := middleware.RoleFromContext(ctx)
	if okay != true{
		log.Println("Error Retrieving Role from Context", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := Handles.RejectStudent(tx, ctx, UserId, Role)
	if err != nil {
		log.Println("Error rejecting applicants from database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit(ctx)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Println("Error encording applicants to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
