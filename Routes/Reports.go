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

	// Get report data
	reportData, total, err := Handles.GetComprehensiveReport(Db.DB, ctx, limit, offset)
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
