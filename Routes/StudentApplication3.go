package Routes

import (
	"log"
	"net/http"
	"time"
	"ubaps/Db"
	"ubaps/Handles"
	middleware "ubaps/Middleware"
	"ubaps/utils"
)

func SubmitForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	Ctx := r.Context() // ✅ FIX

	tx, err := Db.DB.Begin(Ctx)
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(Ctx)

	userId, ok := middleware.UserIDFromContext(Ctx)
	log.Println(userId)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	log.Println("User ID:", userId)

	email, err := Handles.GetEmailByUserID(userId, tx)
	if err != nil {
		http.Error(w, "Failed to get email", http.StatusInternalServerError)
		return
	}

	RegNumber := Handles.GetRegNumberFromEmail(email)

	dob, err := Handles.TimePtr(r.FormValue("dob"))
	if err != nil {
		http.Error(w, "Invalid date of birth", http.StatusBadRequest)
		return
	}

	income, err := Handles.FloatPtr(r.FormValue("monthlyIncome"))
	if err != nil {
		http.Error(w, "Invalid income", http.StatusBadRequest)
		return
	}

	submission := time.Now()

	err = utils.UpdateApplication(
		Ctx,
		tx,
		"submitted",
		dob,
		Handles.StrPtr(r.FormValue("gender")),
		Handles.StrPtr(r.FormValue("Home District")),
		Handles.StrPtr("Computer Engineering"),
		Handles.StrPtr(RegNumber),
		Handles.StrPtr(r.FormValue("Type of intake")),
		Handles.StrPtr(r.FormValue("Accomodation")),
		Handles.StrPtr(r.FormValue("Gurdian Status")),
		Handles.StrPtr(r.FormValue("Guardian Employment Status")),
		income,
		Handles.StrPtr(r.FormValue("otherSupport")),
		Handles.StrPtr(r.FormValue("Reason")),
		&submission,
		userId, // ✅ already int64
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(Ctx); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("You have successfully submitted the application form"))
}
