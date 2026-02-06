package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
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
	results,err := Handles.Applicants(Db.DB,ctx,Pplicants)
	if err != nil{
		log.Println("Error retriving applicants from database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//log.Println("Received applicants data:", Pplicants)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)
	if err != nil{
		log.Println("Error encording applicants to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}