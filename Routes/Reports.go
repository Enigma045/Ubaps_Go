package Routes

import (
	"encoding/json"
	"log"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
	"ubaps/utils"
)

// ExportComprehensiveReport exports the report as JSON (paginated)
func ExportComprehensiveReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	page, limit, offset := utils.GetPaginationParams(r)
	year := r.URL.Query().Get("year")
	semester := r.URL.Query().Get("semester")
	month := r.URL.Query().Get("month")

	var filters Handles.ReportFilters
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&filters)
		if err != nil {
			log.Println("Error decoding report filters:", err)
			// Continue with empty filters if decoding fails, or we could return an error
		}
	}

	// Get report data
	reportData, total, err := Handles.GetComprehensiveReport(Db.DB, ctx, limit, offset, year, semester, month, filters)
	if err != nil {
		log.Println("Reports Handler Error:", err)
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(utils.PaginatedResponse{
		Data:  reportData,
		Total: total,
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		log.Println("JSON encoding error:", err)
	}
}
