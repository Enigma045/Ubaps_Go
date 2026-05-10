package Routes

import (
	"encoding/json"
	"net/http"
	"ubaps/Db"
	"ubaps/Handles"
)

func FinancialReportsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reportType := r.URL.Query().Get("type")

	var result Handles.FinanceReportResult
	var err error

	switch reportType {
	case "disbursement":
		result, err = Handles.GetDisbursementStatusReport(Db.DB, ctx)
	case "utilisation":
		result, err = Handles.GetFundUtilisationReport(Db.DB, ctx)
	case "history":
		result, err = Handles.GetPaymentHistoryReport(Db.DB, ctx)
	case "pending":
		result, err = Handles.GetPendingDisbursementsReport(Db.DB, ctx)
	default:
		http.Error(w, "Invalid report type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
