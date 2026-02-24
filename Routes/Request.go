package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"ubaps/Db"
	"ubaps/Handles"
	"ubaps/utils"
)

func GetRequest_Info(w http.ResponseWriter,r *http.Request) {

    
	
	var request []string
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	results,err := Handles.GetRequest_Info(Db.DB,ctx,request)
	if err != nil{
		log.Println("Error retriving financial Requests from database", err)
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


func AcceptRequest(w http.ResponseWriter, r *http.Request) {

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
	//
    //email := fmt.Sprintf("%s@unilia.ac.mw",student.StudentID)

	//UserId, err := Handles.GetUserIDByEmail(email,tx)
	//if err != nil {
	//	log.Println("Error getting user ID from email:", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
    //

    reqid := strings.Split(student.StudentID, "#")
	log.Println(reqid)
	id,err := utils.Strtoint64(reqid[0])
    if err != nil {
		log.Println("Error Converting String to int64", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := Handles.AcceptRequest(tx, ctx, id)
	if err != nil {
		log.Println("Error Accept Request from database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit(ctx)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Println("Error encording Request to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func RejectRequest(w http.ResponseWriter, r *http.Request) {
	
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
    //email := fmt.Sprintf("%s@unilia.ac.mw",student.StudentID)

	//UserId, err := Handles.GetUserIDByEmail(email,tx)
	//if err != nil {
	//	log.Println("Error getting user ID from email:", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	reqid := strings.Split(student.StudentID, "#")
	log.Println(reqid)
	id,err := utils.Strtoint64(reqid[0])
	if err != nil {
		log.Println("Error Converting String to int64", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := Handles.RejectRequest(tx, ctx, id)
	if err != nil {
		log.Println("Error rejecting request from database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit(ctx)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Println("Error encording request to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func GetTotalAmount(w http.ResponseWriter,r *http.Request){
	ctx := r.Context()



	results,err := Handles.GetTotalAmount(Db.DB,ctx)
	if err != nil{
		log.Println("Error Retrieving Total Amount From Database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)
	if err != nil{
		log.Println("Error encording Total Amount to frontend", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}