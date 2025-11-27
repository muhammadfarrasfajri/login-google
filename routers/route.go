package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/controllers"
)

func SetupRoutes(r *gin.Engine) {
	// Login dengan Google
	r.POST("/api/login", controllers.GoogleAuth)

	// Register user Google
	r.POST("/api/register", controllers.RegisterGoogle)
}
