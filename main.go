package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/config"
	"github.com/muhammadfarrasfajri/login-google/controllers"
	"github.com/muhammadfarrasfajri/login-google/database"
	"github.com/muhammadfarrasfajri/login-google/repository"
	routes "github.com/muhammadfarrasfajri/login-google/routers"
	"github.com/muhammadfarrasfajri/login-google/services"
)

func main() {

	// Init Firebase & MySQL
	config.InitFirebase()
	database.ConnectMySQL()

	app, _ := config.FirebaseApp.Auth(context.Background())

	// Repository
	userRepo := &repository.UserRepository{
		DB: database.DB,
	}

	// Services
	authService := &services.AuthService{
		UserRepo:     userRepo,
		FirebaseAuth: app,
	}

	// Controller
	authController := &controllers.AuthController{
		AuthService: authService,
	}

	// userController := &controllers.UserController{
	// 	UserRepo: userRepo,
	// }

	// Init GIN
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ROUTES
	routes.SetupRoutes(r, authController, userController)

	// Run server
	r.Run(":8080")
}
