package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	user_logs "ubaps/Audit_logs"
	"ubaps/Db"
	"ubaps/utils"
)

func UserLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page, limit, offset := utils.GetPaginationParams(r)

	logs, total, err := user_logs.Get_User_Logs(Db.DB, ctx, limit, offset)
	if err != nil {
		log.Print("Error fetching user logs:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.PaginatedResponse{
		Data:  logs,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func PaymentLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page, limit, offset := utils.GetPaginationParams(r)

	logs, total, err := user_logs.Get_Payment_Logs(Db.DB, ctx, limit, offset)
	if err != nil {
		log.Print("Error fetching payment logs:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.PaginatedResponse{
		Data:  logs,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func ApplicationLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page, limit, offset := utils.GetPaginationParams(r)

	logs, total, err := user_logs.Get_Application_Logs(Db.DB, ctx, limit, offset)
	if err != nil {
		log.Print("Error fetching application logs:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.PaginatedResponse{
		Data:  logs,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}