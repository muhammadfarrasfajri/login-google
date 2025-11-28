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
	config.InitFirebase()
	database.ConnectMySQL()

	app, _ := config.FirebaseApp.Auth(context.Background())

	userRepo := &repository.UserRepository{DB: database.DB}
	authService := &services.AuthService{UserRepo: userRepo, FirebaseAuth: app}
	authController := &controllers.AuthController{AuthService: authService}

	r := gin.Default()
	
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

	routes.SetupRoutes(r, authController)

	r.Run(":8080")
}
