package Routes

import (
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
	middleware "ubaps/Middleware"
	"ubaps/utils"
)	

func Approval(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx := r.Context()
	tx, err := Db.DB.Begin(ctx)

	if err != nil {
		http.Error(w, "Database Failed transctions Error ", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)
    formdata,err := utils.Formdata(r)
	if err != nil {
	   log.Println("Failed to retrieve formdata")
       http.Error(w, "Formdate Error ", http.StatusInternalServerError)
	   return
	}

	userid,check := middleware.UserIDFromContext(ctx)
	if check != true{
	   log.Println("Failed to retrieve userid")
       http.Error(w, "user id ", http.StatusInternalServerError)
	   return
	}

	err = Handles.Db_operation(tx,ctx,userid,0,formdata)
    if err != nil {
	http.Error(w, "Database Failed transctions Error ", http.StatusInternalServerError)
	return
    }
	err = tx.Commit(ctx)
	if err != nil {
		http.Error(w, "Database Failed to commit Error ", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Fees statement successufuly sent"))
}
