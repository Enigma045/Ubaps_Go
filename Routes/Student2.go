package Routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

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
