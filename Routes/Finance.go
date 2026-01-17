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
