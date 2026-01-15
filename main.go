package main

import (
	"log"
	"net/http"

	"ubaps/Db"
	middleware "ubaps/Middleware"
	"ubaps/Routes"
)

func main() {
	// Connect to DB
	Db.ConnectDB()
	defer Db.DB.Close()

	mux := http.NewServeMux()

	/*
		|--------------------------------------------------------------------------
		| API endpoints
		|--------------------------------------------------------------------------
	*/
	mux.HandleFunc("/register", Routes.Register)
	mux.HandleFunc("/verify-email", Routes.VerifyEmail)
	mux.HandleFunc("/Authorize", Routes.Login)
	mux.Handle("/SubmitForm", middleware.RequireAuth(http.HandlerFunc(Routes.SubmitForm)))
	/*
		|--------------------------------------------------------------------------
		| Public static assets (CSS, JS, images)
		|--------------------------------------------------------------------------
	*/
	// Serve CSS/JS with correct MIME type using FileServer
	mux.Handle("/Css/", http.StripPrefix("/Css/", http.FileServer(http.Dir("Pages/Html/student/Css"))))
	mux.Handle("/Js/", http.StripPrefix("/Js/", http.FileServer(http.Dir("Pages/Html/student/Js"))))
	mux.Handle("/Image/", http.StripPrefix("/Image/", http.FileServer(http.Dir("Pages/Html/student/Image"))))
	/*
		|--------------------------------------------------------------------------
		| Public HTML pages (login, register)
		|--------------------------------------------------------------------------
	*/
	mux.Handle("/Html/", http.StripPrefix("/Html/", http.FileServer(http.Dir("Pages/student/public"))))

	/*
		|--------------------------------------------------------------------------
		| Login page route
		|--------------------------------------------------------------------------
	*/
	mux.HandleFunc("/Login", Routes.Login_page)
	mux.HandleFunc("/", Routes.Sign_Up_page)
	/*
		|--------------------------------------------------------------------------
		| Logout route (authenticated)
		|--------------------------------------------------------------------------
	*/
	mux.Handle(
		"/logout",
		middleware.RequireAuth(
			http.HandlerFunc(middleware.Logout),
		),
	)

	/*
		|--------------------------------------------------------------------------
		| Protected student dashboard route
		|--------------------------------------------------------------------------
	*/
	mux.Handle(
		"/dashboard",
		middleware.RequireAuth(
			middleware.RequireRole("student")(
				http.HandlerFunc(Routes.StudentDashboard),
			),
		),
	)

	mux.Handle(
		"/apply",
		middleware.RequireAuth(
			middleware.RequireRole("student")(
				http.HandlerFunc(Routes.ApplicationForm),
			),
		),
	)

	/*
		|--------------------------------------------------------------------------
		| Start server
		|--------------------------------------------------------------------------
	*/
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// http.Handle(
// 	"/admin/dashboard",
// 	middleware.RequireAuth(
// 		middleware.RequireRole("admin")(
// 			http.HandlerFunc(AdminDashboard),
// 		),
// 	),
// )

// http.Handle(
// 	"/logout",
// 	middleware.RequireAuth(
// 		http.HandlerFunc(middleware.Logout),
// 	),
// )
// mux.Handle(
// 		"/dashboard",
// 		middleware.RequireAuth(
// 			http.HandlerFunc(Routes.Dashboard), // example
// 		),
// 	)
