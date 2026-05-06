package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
	middleware "ubaps/Middleware"
)

func StatsStudent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	stats, err := Handles.GetStudentStats(Db.DB, ctx, userID)
	if err != nil {
		log.Println("Error fetching student stats:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

    log.Println(stats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func StatsRegistrar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := Handles.GetRegistrarStats(Db.DB, ctx)
	if err != nil {
		log.Println("Error fetching registrar stats:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println(stats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func StatsAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := Handles.GetAdminStats(Db.DB, ctx)
	if err != nil {
		log.Println("Error fetching admin stats:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

    log.Println(stats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func StatsDean(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := Handles.GetDeanStats(Db.DB, ctx)
	if err != nil {
		log.Println("Error fetching dean stats:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println(stats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func StatsFinance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := Handles.GetFinanceStats(Db.DB, ctx)
	if err != nil {
		log.Println("Error fetching finance stats:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println(stats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func UserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := Handles.GetUserProfile(Db.DB, ctx, userID)
	if err != nil {
		log.Println("Error fetching user profile:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println(profile)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func GetDetailedProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := Handles.GetDetailedUserProfile(Db.DB, ctx, userID)
	if err != nil {
		log.Println("Error fetching detailed user profile:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
