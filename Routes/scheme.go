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



func Scheme_Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx := r.Context()
	// Role, rcheck := middleware.RoleFromContext(ctx)
	// if rcheck != true {
	// 	log.Println("Failed to take UserId")
	// }
	// Begin transaction
	tx, err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx) // Rollback will be ignored if commit succeeds
	//
	userId, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	log.Println("User ID:", userId)
	// Parse form safely
	formData, err := utils.Formdata(r)
	if err != nil {
		log.Println("Formdata Error:", err)
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	log.Println("FormData received:", formData)

	// Execute DB operation safely
	if err := utils.Scheme_Operations(tx, ctx, formData, userId); err != nil {
		log.Println("DB Operation Failed:", err)
		http.Error(w, "Database operation failed", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You have successfully submitted the application form"))
}

func GetBenefactor(w http.ResponseWriter, r *http.Request){
   ctx := r.Context()
   w.Header().Set("Content-Type", "application/json")
   //userId, ok := middleware.UserIDFromContext(ctx)
   //if !ok {
	//	log.Println("Failed to take UserId")
	//	return
//}

	benefactor, err := utils.GetBenefactor(Db.DB, ctx)
	if err != nil {
		log.Println("Failed to get benefactor:", err)
		http.Error(w, "Failed to get benefactor", http.StatusInternalServerError)
		return
	}

	log.Println("Benefactor:", benefactor)


	err = json.NewEncoder(w).Encode(benefactor)
	if err != nil {
		log.Println("Failed to encode benefactor:", err)
		http.Error(w, "Failed to encode benefactor", http.StatusInternalServerError)
		return
	}

	
}


func DeleteBenefactor(w http.ResponseWriter,r *http.Request){
    var emailreq emailRequest

	if r.Method != http.MethodPost{
    http.Error(w,"wrong Method",http.StatusMethodNotAllowed)
	}

	ctx := r.Context()

	tx,err := Db.DB.Begin(ctx)
    if err != nil {
		log.Println(err)
		http.Error(w,"Failed to create Transction",http.StatusInternalServerError)
	    return
	}
	defer tx.Rollback(ctx)

	err = json.NewDecoder(r.Body).Decode(&emailreq)
	if err != nil {
		log.Println(err)
		http.Error(w,"Invalid JSON",http.StatusBadRequest)
		return
	}
	
	// userid,err := Handles.GetUserIDByEmail(emailreq.Email,tx)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w,"Failed to get userid",http.StatusInternalServerError)
	//     return
	// }

	var name string = emailreq.Name

	err = utils.DeleteBenefactor(tx,ctx,name)
	if err != nil {
		log.Println(err)
		http.Error(w,"Invalid JSON",http.StatusBadRequest)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit failed:", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You have succefully deleted the user"))

}

func GetScheme(w http.ResponseWriter,r *http.Request){


	w.Header().Set("Content-Type", "application/json")
    ctx := r.Context()

	schemes,err := utils.GetScheme(Db.DB,ctx)
    if err != nil {
		log.Println(err)
		http.Error(w,"Failed to get schemes",http.StatusInternalServerError)
	    return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schemes)
}

type SchemeInfo struct {
	Reg string `json:"reg"`
	Scheme string `json:"scheme"`
	Amount string `json:"amount"`
}

func SendScheme_Info(w http.ResponseWriter,r *http.Request){


	w.Header().Set("Content-Type", "application/json")
    ctx := r.Context()
	tx,err := Db.DB.Begin(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w,"Failed to create Transction",http.StatusInternalServerError)
	    return
	}
	defer tx.Rollback(ctx)

    var schemeinfo SchemeInfo

    err = json.NewDecoder(r.Body).Decode(&schemeinfo)
    if err != nil {
		log.Println(err)
		http.Error(w,"Invalid JSON",http.StatusBadRequest)
		return
	}

	email := fmt.Sprintf("%s@unilia.ac.mw",schemeinfo.Reg)
	userid,err := Handles.GetUserIDByEmail(email,tx)
    if err != nil {
		log.Println(err)
		http.Error(w,"Failed to get userid",http.StatusInternalServerError)
	    return
	}

	exist,err :=utils.CheckForScheme(tx,ctx,userid)
	if err != nil {
		log.Println(err)
		http.Error(w,"Failed to check for scheme",http.StatusInternalServerError)
	    return
	}

	if exist {
		log.Println("User already has a scheme")
		http.Error(w,"User already has a scheme",http.StatusBadRequest)
	    return
	}

    schemeid,err := utils.GetSchemeId(schemeinfo.Scheme,tx,ctx)
	if err != nil {
		log.Println(err)
		http.Error(w,"Failed to get schemeid",http.StatusInternalServerError)
	    return
	}

	err = utils.CheckSchemeAmount(schemeinfo.Scheme,tx,ctx,schemeinfo.Amount)
	if err != nil {
		log.Println(err)
		http.Error(w,"Amount is less then bursary scheme amount",http.StatusInternalServerError)
	    return
	}

	value,err := utils.GetAvailableAmount(tx,ctx,schemeid)
    if err != nil{
		log.Println(err)
		http.Error(w,"Failed to retrieve availabele scheme balnce",http.StatusInternalServerError)
		return 
	}

	err = utils.UpdateScheme_Amount(tx,ctx,schemeid,schemeinfo.Amount,value)
    if err != nil{
		log.Println(err)
		http.Error(w,"Failed to update scheme amount scheme balnce",http.StatusInternalServerError)
		return 
	}

	err = utils.SendScheme_Info(tx,ctx,userid,schemeid,schemeinfo.Amount)
	if err != nil {
		log.Println(err)
		http.Error(w,"Failed to send scheme info",http.StatusInternalServerError)
		return
	}

	tx.Commit(ctx)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You have succefully sent the scheme info"))
}
