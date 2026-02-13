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
	mux.HandleFunc("/fees", Routes.Fees)
	mux.Handle("/sendrequest", middleware.RequireAuth(http.HandlerFunc(Routes.Approval)))
	mux.Handle("/getnotifications", middleware.RequireAuth(http.HandlerFunc(Routes.Notifications)))
	mux.Handle("/countnotifications", middleware.RequireAuth(http.HandlerFunc(Routes.NotificationCounter)))
	mux.Handle("/benefactor", middleware.RequireAuth(http.HandlerFunc(Routes.Scheme_Info)))
	mux.Handle("/SubmitForm", middleware.RequireAuth(http.HandlerFunc(Routes.SubmitForm)))
	//admin routes
	mux.Handle("/getuserdetails", middleware.RequireAuth(http.HandlerFunc(Routes.Getuserdetails)))
	mux.Handle("/deleteaccount", middleware.RequireAuth(http.HandlerFunc(Routes.DeleteAccount)))
	mux.Handle("/createuser", middleware.RequireAuth(http.HandlerFunc(Routes.CreateUser)))
	mux.Handle("/userlog", middleware.RequireAuth(http.HandlerFunc(Routes.UserLogs)))
	mux.Handle("/paymentlog", middleware.RequireAuth(http.HandlerFunc(Routes.PaymentLogs)))
	mux.Handle("/applicationlog", middleware.RequireAuth(http.HandlerFunc(Routes.ApplicationLogs)))
	//registrar routes
	mux.Handle("/applicants", middleware.RequireAuth(http.HandlerFunc(Routes.Applicants)))
	mux.Handle("/getbenefactor", middleware.RequireAuth(http.HandlerFunc(Routes.GetBenefactor)))
	mux.Handle("/deletebenefactor", middleware.RequireAuth(http.HandlerFunc(Routes.DeleteBenefactor)))
	mux.Handle("/considerstudent", middleware.RequireAuth(http.HandlerFunc(Routes.ConsiderStudent)))
	mux.Handle("/rejectstudent", middleware.RequireAuth(http.HandlerFunc(Routes.RejectStudent)))
	mux.Handle("/getschemes", middleware.RequireAuth(http.HandlerFunc(Routes.GetScheme)))
	mux.Handle("/schemeinfo", middleware.RequireAuth(http.HandlerFunc(Routes.SendScheme_Info)))
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
	mux.HandleFunc("/commitee", Routes.Commitee)
	mux.HandleFunc("/Login", Routes.Login_page)
	mux.HandleFunc("/", Routes.Sign_Up_page)
	mux.HandleFunc("/request", Routes.Request_Page)
	//mux.HandleFunc("/financial", Routes.Approval_Page)

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

	mux.Handle(
		"/notifications",
		middleware.RequireAuth(
			middleware.RequireRole("student")(
				http.HandlerFunc(Routes.Notification_Page),
			),
		),
	)

	mux.Handle(
		"/letters",
		middleware.RequireAuth(
			middleware.RequireRole("student")(
				http.HandlerFunc(Routes.Letter_Page),
			),
		),
	)
	/*
		|--------------------------------------------------------------------------
		| Protected financial office dashboard route
		|--------------------------------------------------------------------------
	*/
	mux.Handle(
		"/financialdashboard",
		middleware.RequireAuth(
			middleware.RequireRole("finance_office")(
				http.HandlerFunc(Routes.FinancialDashboard_Page),
			),
		),
	)
	mux.Handle(
		"/financial",
		middleware.RequireAuth(
			middleware.RequireRole("finance_office")(
				http.HandlerFunc(Routes.Approval_Page),
			),
		),
	)

	mux.Handle(
		"/financial_request",
		middleware.RequireAuth(
			middleware.RequireRole("finance_office")(
				http.HandlerFunc(Routes.Request_Page),
			),
		),
	)
	mux.Handle(
		"/financial_reports",
		middleware.RequireAuth(
			middleware.RequireRole("finance_office")(
				http.HandlerFunc(Routes.Financial_reports_Page),
			),
		),
	)
	/*
		|--------------------------------------------------------------------------
		| Protected admin dashboard route
		|--------------------------------------------------------------------------
	*/
	mux.Handle(
		"/admindashboard",
		middleware.RequireAuth(
			middleware.RequireRole("admin")(
				http.HandlerFunc(Routes.AdminDashboard_Page),
			),
		),
	)

	mux.Handle(
		"/adminusers",
		middleware.RequireAuth(
			middleware.RequireRole("admin")(
				http.HandlerFunc(Routes.AdminUsers_Page),
			),
		),
	)

	mux.Handle(
		"/admintrails",
		middleware.RequireAuth(
			middleware.RequireRole("admin")(
				http.HandlerFunc(Routes.AdminTrails_Page),
			),
		),
	)

	mux.Handle(
		"/adminmodification",
		middleware.RequireAuth(
			middleware.RequireRole("admin")(
				http.HandlerFunc(Routes.AdminModification_Page),
			),
		),
	)
	/*
		|--------------------------------------------------------------------------
		| Protected dean of students dashboard route
		|--------------------------------------------------------------------------
	*/
	mux.Handle(
		"/deandashboard",
		middleware.RequireAuth(
			middleware.RequireRole("dean_of_student")(
				http.HandlerFunc(Routes.DeanDashboard_Page),
			),
		),
	)

	/*
		|--------------------------------------------------------------------------
		| Protected registrar dashboard route
		|--------------------------------------------------------------------------
	*/
	mux.Handle(
		"/registrardashboard",
		middleware.RequireAuth(
			middleware.RequireRole("registrar")(
				http.HandlerFunc(Routes.RegistrarDashboard_Page),
			),
		),
	)
	mux.Handle(
		"/decision",
		middleware.RequireAuth(
			middleware.RequireRole("registrar")(
				http.HandlerFunc(Routes.Decision_Page),
			),
		),
	)

	mux.Handle(
		"/scheme",
		middleware.RequireAuth(
			middleware.RequireRole("registrar")(
				http.HandlerFunc(Routes.Scheme_page),
			),
		),
	)
	mux.Handle(
		"/receive_letters",
		middleware.RequireAuth(
			middleware.RequireRole("registrar")(
				http.HandlerFunc(Routes.Letters_page),
			),
		),
	)
	mux.Handle(
		"/financial_approval",
		middleware.RequireAuth(
			middleware.RequireRole("registrar")(
				http.HandlerFunc(Routes.Financial_approval_page),
			),
		),
	)
	mux.Handle(
		"/general_reports",
		middleware.RequireAuth(
			middleware.RequireRole("registrar")(
				http.HandlerFunc(Routes.General_reports_Page),
			),
		),
	)
	// mux.Handle(
	// 	"/scheme",
	// 	middleware.RequireAuth(
	// 		middleware.RequireRole("registrar")(
	// 			http.HandlerFunc(Routes.Scheme_page),
	// 		),
	// 	),
	// )
	/*
		|--------------------------------------------------------------------------
		| Logout route (authenticated)
		|--------------------------------------------------------------------------
	*/
	mux.HandleFunc("/logout", middleware.Logout)

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
