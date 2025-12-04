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
	adminRepo := &repository.AdminRepository{
		DB: database.DB,
	}
	userRepo := &repository.UserRepository{
		DB: database.DB,
	}

	jwtManager := &middleware.JWTManager{
		 AccessSecret:  []byte(os.Getenv("JWT_SECRET")),
		 RefreshSecret: []byte(os.Getenv("REFRESH_SECRET")),
		}

	// Services
	authAdminService := &services.AuthService{
		Repo: adminRepo,
		FirebaseAuth: appAdmin,
		JWTSecret: jwtManager,
	}

	authUserService := &services.AuthService{
		Repo: userRepo,
		FirebaseAuth: appUser,
		JWTSecret: jwtManager,
	}

	userService := &services.UserService{
		UserRepo: userRepo,
	}
	
	// Controller
	authAdminController := &controllers.AuthController{
		AuthService: authAdminService,
	}
	authUserController := &controllers.AuthController{
		AuthService: authUserService,
	}

	userController := &controllers.UserController{
		UserService: userService,
		Repo:        userRepo,
	}

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
