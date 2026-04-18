package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
)

func GetStatementRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	results, err := Handles.GetStatementRequests(Db.DB, ctx)
	if err != nil {
		log.Println("Error retrieving statement requests from database:", err)
		http.Error(w, "Failed to retrieve statement requests", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
    // Results is an array of StatementRequest
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Println("Error encoding statement requests:", err)
		http.Error(w, "Error returning data", http.StatusInternalServerError)
		return
	}
}
