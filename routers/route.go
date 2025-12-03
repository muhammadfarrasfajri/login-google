package routes

import (
	"net/http"

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
			c.JSON(http.StatusOK, gin.H{
				"message": "logout success",
			})
		})
	}

	user := r.Group("/users", middleware.AuthMiddleware())
	{
		user.GET("/:id", userController.GetByID)
	}
	
	// ===========================
	// ADMIN ROUTES
	// ===========================
	admin := r.Group("/admin")
	{
		admin.GET("/users", userController.GetAll)
		admin.GET("/:id", userController.GetByID)
		admin.PATCH("/users/:id", userController.Update)
		admin.DELETE("/users/:id", userController.Delete)
	}
}