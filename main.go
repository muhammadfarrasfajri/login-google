package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/muhammadfarrasfajri/login-google/config"
	"github.com/muhammadfarrasfajri/login-google/controllers"
	"github.com/muhammadfarrasfajri/login-google/database"
	"github.com/muhammadfarrasfajri/login-google/middleware"
	"github.com/muhammadfarrasfajri/login-google/repository"
	routes "github.com/muhammadfarrasfajri/login-google/routers"
	"github.com/muhammadfarrasfajri/login-google/services"
)

func main() {

	// Init Firebase & MySQL
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found!")
	}
	if os.Getenv("JWT_SECRET") == "" || os.Getenv("REFRESH_SECRET") == "" {
		log.Fatal("JWT secrets must not be empty!")
	}
	
	middleware.InitEncryptionKey()
	config.InitFirebase()
	database.ConnectMySQL()

	appAdmin, _ := config.FirebaseAppAdmin.Auth(context.Background())
	appUser, _ := config.FirebaseAppUser.Auth(context.Background())

	// Repository
	adminRepo := repository.NewAdminRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)

	// JWT
	jwtManager := middleware.NewJWTManager(os.Getenv("JWT_SECRET"),os.Getenv("REFRESH_SECRET"))

	// Services
	//Auth Admin and User
	authAdminService := services.NewAuthService(adminRepo, appAdmin, jwtManager)
	authUserService := services.NewAuthService(userRepo, appUser, jwtManager)
	//CRUD User
	userService := services.NewUserSevice(userRepo)
	
	// Controller
	// Auth Admin and User
	authAdminController := controllers.NewAuthController(authAdminService)
	authUserController := controllers.NewAuthController(authUserService)
	//CRUD User
	userController := controllers.NewUserController(userService, userRepo)

	// Init GIN
	r := gin.Default()

	r.Static("/public", "./public")

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ROUTES
	routes.SetupRoutes(r, authAdminController, authUserController, userController, jwtManager)

	// Run server
	r.Run(":8080")
}
