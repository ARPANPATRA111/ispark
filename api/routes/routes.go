package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/controllers"
	"github.com/iips-oss/ispark/api/middleware"
)

// SetupRoutes configures the endpoints for the API
func SetupRoutes(app *fiber.App) {
	// Base API group
	api := app.Group("/api")

	// Auth group
	auth := api.Group("/auth")

	// Public Auth routes
	auth.Get("/captcha", controllers.GetCaptcha)
	auth.Post("/register", controllers.Register)
	auth.Post("/verify-otp", controllers.VerifyOTP)
	auth.Post("/login", controllers.Login)
	auth.Post("/forgot-password", controllers.ForgotPassword)
	auth.Post("/reset-password", controllers.ResetPassword)
	auth.Post("/refresh", controllers.RefreshToken)

	// Protected routes (Require login)
	auth.Use(middleware.AuthRequired())
	auth.Post("/logout", controllers.Logout)
	auth.Get("/profile", controllers.GetProfile)

	// Student Dashboard routes (Require login)
	student := api.Group("/student")
	student.Use(middleware.AuthRequired())
	student.Get("/certificates", controllers.GetCertificates)
	student.Post("/certificates", controllers.UploadCertificate)
	student.Get("/certificates/:id/file", controllers.DownloadCertificate)
	student.Get("/leaderboard/champions", controllers.GetCategoryChampions)
	student.Get("/leaderboard", controllers.GetLeaderboard)
	student.Get("/activities", controllers.GetActivities)
	student.Put("/profile", controllers.UpdateProfile)
	student.Post("/change-password", controllers.ChangePassword)
	student.Post("/activities/:id/enroll", controllers.EnrollActivity)
	student.Get("/enrollments", controllers.GetEnrollments)
	student.Get("/dashboard/stats", controllers.GetDashboardStats)
	student.Get("/marksheet", controllers.GetMarksheet)

	// Admin
	api.Post("/admin/auth/login", controllers.AdminLogin)

	admin := api.Group("/admin")

	// Must be logged in AND hold an administrative role. A super admin is not
	// batch-scoped, so it sees every student.
	admin.Use(middleware.AuthRequired())
	admin.Use(middleware.RoleRequired("admin", "superadmin"))

	// Must change the password
	admin.Post("/change-password", controllers.AdminChangePassword)

	admin.Get("/students", controllers.GetAllStudents)
	admin.Get("/students/:roll", controllers.GetStudentDetail)

	// Platform-wide routes, super admin only
	platform := admin.Group("/platform", middleware.RoleRequired("superadmin"))
	platform.Get("/stats", controllers.GetPlatformStats)
	platform.Get("/users", controllers.GetPlatformUsers)
	platform.Post("/users", controllers.CreatePlatformUser)
	platform.Delete("/users/:id", controllers.DeletePlatformUser)

}
