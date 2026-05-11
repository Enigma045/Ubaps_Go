package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	user_logs "ubaps/Audit_logs"
	"ubaps/Db"
	"ubaps/utils"
)

func parseLogFilters(r *http.Request) user_logs.LogFilters {
	filters := user_logs.LogFilters{}

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			filters.UserID = &id
		}
	}
	if targetUserIDStr := r.URL.Query().Get("target_user_id"); targetUserIDStr != "" {
		if id, err := strconv.ParseInt(targetUserIDStr, 10, 64); err == nil {
			filters.TargetUserID = &id
		}
	}
	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = &t
		}
	}
	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set to end of day
			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filters.EndDate = &t
		}
	}

	return filters
}

func UserLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page, limit, offset := utils.GetPaginationParams(r)
	filters := parseLogFilters(r)

	logs, total, err := user_logs.Get_User_Logs(Db.DB, ctx, limit, offset, filters)
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
	filters := parseLogFilters(r)

	logs, total, err := user_logs.Get_Payment_Logs(Db.DB, ctx, limit, offset, filters)
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
	filters := parseLogFilters(r)

	logs, total, err := user_logs.Get_Application_Logs(Db.DB, ctx, limit, offset, filters)
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

func AllLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page, limit, offset := utils.GetPaginationParams(r)
	filters := parseLogFilters(r)

	typesStr := r.URL.Query().Get("types")
	var types []string
	if typesStr != "" {
		types = strings.Split(typesStr, ",")
	}

	logs, total, err := user_logs.Get_Unified_Logs(Db.DB, ctx, limit, offset, types, filters)
	if err != nil {
		log.Print("Error fetching unified logs:", err)
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