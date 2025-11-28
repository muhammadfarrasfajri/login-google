package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/controllers"
	"github.com/muhammadfarrasfajri/login-google/middleware"
)

func SetupRoutes(r *gin.Engine, authController *controllers.AuthController, userController *controllers.UserController) {

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
	// USER ROUTES (sementara NO AUTH untuk testing)
	// ===========================
	user := r.Group("/users", middleware.AuthMiddleware()) // enable kalo sudah siap
	// user := r.Group("/users") // sementara TANPA auth biar testing enak

	// ===========================
	// ADMIN ONLY ROUTES
	// ===========================
	{
		user.GET("/", userController.GetAll)
		user.GET("/:id", userController.GetByID)
		user.PUT("/:id", userController.Update)
		user.DELETE("/:id", userController.Delete)
	}

	// ===========================
	// ADMIN ROUTES
	// ===========================
	admin := r.Group("/admin", middleware.AuthMiddleware(), middleware.AdminOnly())
	// admin := r.Group("/admin") // sementara tanpa auth
	{
		admin.GET("/users", userController.GetAll)
		admin.PUT("/users/:id", userController.Update)
		admin.DELETE("/users/:id", userController.Delete)
	}
}