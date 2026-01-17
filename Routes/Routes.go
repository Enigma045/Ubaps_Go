package Routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Request_Page(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("Pages/Html/student/public/Finance_Request.html")
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
