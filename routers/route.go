package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/controllers"
	"github.com/muhammadfarrasfajri/login-google/middleware"
)

func SetupRoutes(r *gin.Engine, authAdminController *controllers.AuthController, authUserController *controllers.AuthController, userController *controllers.UserController, paymentController *controllers.PaymentController, jwtManager *middleware.JWTManager) {

	// ===========================
	// AUTH ROUTES
	// ===========================
	auth := r.Group("/api/auth")
	{
		//auth admin
		auth.POST("/admin/register", authAdminController.RegisterAdmin)
		auth.POST("/admin/login", authAdminController.LoginAdmin)
		auth.POST("/admin/refresh", authAdminController.RefreshTokenAdmin)
		auth.POST("/admin/logout", jwtManager.AuthMiddleware(), authAdminController.LogoutAdmin)

		//auth user
		auth.POST("/user/register", authUserController.RegisterUser)
		auth.POST("/user/login", authUserController.LoginUser)
		auth.POST("/user/refresh", authUserController.RefreshTokenAdmin)
		auth.POST("/user/logout", jwtManager.AuthMiddleware(), authUserController.LogoutUser)
	}

	// ===========================
	// USER ROUTES
	// ===========================
	user := r.Group("/users", jwtManager.AuthMiddleware())
	{
		user.GET("/:id", userController.GetByID)
	}

	// ===========================
	// ADMIN ROUTES
	// ===========================

	admin := r.Group("/admin", jwtManager.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.GET("/:id", userController.GetByID)
		admin.GET("/users", userController.GetAll)
		admin.PATCH("/users/:id", userController.Update)
		admin.DELETE("/users/:id", userController.Delete)
	}

	// ===========================
	// PAYMENT ROUTES (MIDTRANS)
	// ===========================

	// 1. Webhook / Callback (HARUS PUBLIC)
	// Midtrans akan memanggil URL ini. Jangan gunakan middleware Auth di sini.
	r.POST("/midtrans/callback", paymentController.MidtransCallback)

	// 2. Transaksi QRIS (BUTUH LOGIN)
	// User harus login untuk membuat pembayaran
	payment := r.Group("/api/payment", jwtManager.AuthMiddleware())
	{
		// Endpoint: /api/payment/qris
		payment.POST("/qris", paymentController.CreatePayment)
	}
}
