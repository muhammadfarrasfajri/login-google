package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/controllers"
	"github.com/muhammadfarrasfajri/login-google/middleware"
)

func SetupRoutes(r *gin.Engine, authAdminController *controllers.AuthController,authUserController *controllers.AuthController, userController *controllers.UserController, jwtManager *middleware.JWTManager) {

	// ===========================
	// AUTH ROUTES USERS
	// ===========================
	auth := r.Group("/api/auth")
	{
		auth.POST("/admin/register", authAdminController.RegisterAdmin)
		auth.POST("/admin/login", authAdminController.LoginAdmin)
		auth.POST("/admin/logout",jwtManager.AuthMiddleware(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "logout success",
			})
		})
		auth.POST("/user/register", authUserController.RegisterUser)
		auth.POST("/user/login", authUserController.LoginUser)
		auth.POST("/user/logout", jwtManager.AuthMiddleware(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "logout success",
			})
		})
	}

	user := r.Group("/users", jwtManager.AuthMiddleware())
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