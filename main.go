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
	mux.HandleFunc("/api/forgot-password", Routes.ForgotPasswordAPI)
	mux.HandleFunc("/api/reset-password", Routes.ResetPasswordAPI)
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
	//Financial Office
	mux.Handle("/getrequest", middleware.RequireAuth(http.HandlerFunc(Routes.GetRequest_Info)))
	mux.Handle("/acceptrequest", middleware.RequireAuth(http.HandlerFunc(Routes.AcceptRequest)))
	mux.Handle("/rejectrequest", middleware.RequireAuth(http.HandlerFunc(Routes.RejectRequest)))
	mux.Handle("/gettotalamount", middleware.RequireAuth(http.HandlerFunc(Routes.GetTotalAmount)))
	//Review Card
	mux.Handle("/getapplicationstatus", middleware.RequireAuth(http.HandlerFunc(Routes.GetApplicationStatus)))
	mux.Handle("/requeststatement", middleware.RequireAuth(http.HandlerFunc(Routes.Request_Statement)))
	mux.Handle("/getstatementrequests", middleware.RequireAuth(http.HandlerFunc(Routes.GetStatementRequests)))
	mux.Handle("/api/get-financial-history", middleware.RequireAuth(http.HandlerFunc(Routes.GetFinancialHistory)))

	// Statistics endpoints
	mux.Handle("/stats/student", middleware.RequireAuth(http.HandlerFunc(Routes.StatsStudent)))
	mux.Handle("/stats/registrar", middleware.RequireAuth(http.HandlerFunc(Routes.StatsRegistrar)))
	mux.Handle("/stats/admin", middleware.RequireAuth(http.HandlerFunc(Routes.StatsAdmin)))
	mux.Handle("/stats/dean", middleware.RequireAuth(http.HandlerFunc(Routes.StatsDean)))
	mux.Handle("/stats/finance", middleware.RequireAuth(http.HandlerFunc(Routes.StatsFinance)))
	mux.Handle("/user/profile", middleware.RequireAuth(http.HandlerFunc(Routes.UserProfile)))
	mux.Handle("/user/profile/detailed", middleware.RequireAuth(http.HandlerFunc(Routes.GetDetailedProfile)))
	mux.Handle("/api/user/profile/update", middleware.RequireAuth(http.HandlerFunc(Routes.UpdateProfile)))
	mux.Handle("/api/user/profile/change-password", middleware.RequireAuth(http.HandlerFunc(Routes.ChangePassword)))
	mux.Handle("/api/admin/update-user", middleware.RequireAuth(middleware.RequireRole("admin")(http.HandlerFunc(Routes.AdminUpdateUserAPI))))
	mux.Handle("/api/admin/trigger-reset", middleware.RequireAuth(middleware.RequireRole("admin")(http.HandlerFunc(Routes.TriggerPasswordResetAPI))))

	// Report endpoints
	mux.Handle("/api/reports/comprehensive", middleware.RequireAuth(middleware.RequireRole("registrar", "finance_office", "admin")(http.HandlerFunc(Routes.ExportComprehensiveReport))))
	
	// Letter API endpoints
	mux.Handle("/api/submit-letter", middleware.RequireAuth(middleware.RequireRole("student")(http.HandlerFunc(Routes.SubmitLetter))))
	mux.Handle("/api/get-letters", middleware.RequireAuth(middleware.RequireRole("registrar")(http.HandlerFunc(Routes.GetLettersList))))
	mux.Handle("/api/download-letter", middleware.RequireAuth(middleware.RequireRole("registrar")(http.HandlerFunc(Routes.DownloadLetter))))
	mux.Handle("/api/send-letter", middleware.RequireAuth(middleware.RequireRole("registrar")(http.HandlerFunc(Routes.SendLetter))))

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
	mux.Handle("/payinstallment", middleware.RequireAuth(http.HandlerFunc(Routes.PayInstallment)))
	mux.HandleFunc("/commitee", Routes.Commitee)
	mux.HandleFunc("/Login", Routes.Login_page)
	mux.HandleFunc("/forgot-password", Routes.ForgotPassword_page)
	mux.HandleFunc("/reset-password", Routes.ResetPassword_page)
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
		"/user/profile-settings",
		middleware.RequireAuth(
			http.HandlerFunc(Routes.UserProfileSettings_Page),
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
		"/financial_review",
		middleware.RequireAuth(
			middleware.RequireRole("finance_office")(
				http.HandlerFunc(Routes.Financial_Review_Page),
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
		"/audit_reports",
		middleware.RequireAuth(
			middleware.RequireRole("admin")(
				http.HandlerFunc(Routes.AuditReports_Page),
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
			middleware.RequireRole(
				"dean_of_student",
				"dean_of_facult",
				"dean_of_science",
			)(
				http.HandlerFunc(Routes.DeanDashboard_Page),
			),
		),
	)

	mux.Handle(
		"/deandecision",
		middleware.RequireAuth(
			middleware.RequireRole(
				"dean_of_student",
				"dean_of_facult",
				"dean_of_science",
				"finance_office",
			)(
				http.HandlerFunc(Routes.Dean_Decision_page),
			),
		),
	)
	mux.Handle(
		"/Student_reports",
		middleware.RequireAuth(
			middleware.RequireRole(
				"dean_of_student",
				"dean_of_facult",
				"dean_of_science",
			)(
				http.HandlerFunc(Routes.Student_reports_Page),
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
		"/request_approval",
		middleware.RequireAuth(
			middleware.RequireRole("registrar")(
				http.HandlerFunc(Routes.Financial_request_approval_page),
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
