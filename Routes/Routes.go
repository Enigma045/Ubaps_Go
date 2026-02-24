package Routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Financial_request_approval_page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Financial_Reaquest_approval.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}
func Dean_Decision_page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Deans_Decisions.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}
func Student_reports_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Student_reports.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Letters_page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Letters_receive.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Financial_approval_page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Financial_Approval.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func General_reports_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Combined_reports.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Financial_reports_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Financial_reports.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Decision_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Committe_Review.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func AdminTrails_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Audit_Trailer.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func AdminUsers_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Admin_Users.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func AdminModification_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Admin_Modification.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func AdminDashboard_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Admin_Dashboard.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func FinancialDashboard_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Financial_Dashboard.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func RegistrarDashboard_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Registrar_Dashboard.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func DeanDashboard_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Dean_Dashboard.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Letter_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Student_Portal.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Notification_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Student_Notification.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Approval_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Financial_Approval.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Request_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Finance_Request.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Commitee(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Login_Decision.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}
func Scheme_page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/Protected/Bursary_Scheme.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func StudentDashboard(w http.ResponseWriter, r *http.Request) {

	data, err := os.ReadFile("Pages/Html/student/Protected/Student_Dashboard.html")
	if err != nil {
		log.Println("Page not found")
		http.Error(w, fmt.Sprintf("Page not found %s", err), http.StatusNotFound)
		return
	}
	log.Println("completed")
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func ApplicationForm(w http.ResponseWriter, r *http.Request) {

	data, err := os.ReadFile("Pages/Html/student/Protected/Student_Application.html")
	if err != nil {
		log.Println("Page not found")
		http.Error(w, fmt.Sprintf("Page not found %s", err), http.StatusNotFound)
		return
	}
	log.Println("completed")
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Sign_Up_page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/public/Register.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func Login_page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/public/login.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}
