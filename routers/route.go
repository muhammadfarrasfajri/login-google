package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/controllers"
	"github.com/muhammadfarrasfajri/login-google/middleware"
)

func SetupRoutes(r *gin.Engine, authController *controllers.AuthController) {

	// ===========================
	// AUTH ROUTES
	// ===========================
	auth := r.Group("/api")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)

		auth.POST("/logout", middleware.AuthMiddleware(), func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "logout success",
			})
		})
	}

	// ===========================
	// PROTECTED ROUTES (user harus login)
	// ===========================
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// // USER PROFILE ROUTE
	// protected.GET("/profile", userController.Profile)

	// // ===========================
	// // ADMIN ONLY ROUTES
	// // ===========================
	// admin := protected.Group("/admin")
	// admin.Use(middleware.AdminOnly())
	// {
	// 	admin.GET("/users", userController.GetAllUsers)
	// 	admin.DELETE("/users/:id", userController.DeleteUser)
	// 	admin.PUT("/users/:id", userController.UpdateUser)
	// }
}
