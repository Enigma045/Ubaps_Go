package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
)

type statusRequest struct {
	RegNumber string `json:"reg"`
}

func GetApplicationStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var req statusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := Handles.GetApplicationStatus(Db.DB, ctx, req.RegNumber)
	if err != nil {
		log.Println("Error retrieving application status:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		log.Println("Error encoding application status:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
