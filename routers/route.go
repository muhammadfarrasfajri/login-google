package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/controllers"
)

func SetupRoutes(r *gin.Engine, authController *controllers.AuthController) {
	auth := r.Group("/api")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)

	}
}

