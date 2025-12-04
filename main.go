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

	appAdmin, _ := config.FirebaseAppAdmin.Auth(context.Background())
	appUser, _ := config.FirebaseAppUser.Auth(context.Background())

	// Repository
	adminRepo := &repository.AdminRepository{
		DB: database.DB,
	}
	userRepo := &repository.UserRepository{
		DB: database.DB,
	}

	// Services
	authAdminService := &services.AuthService{
		Repo: adminRepo,
		FirebaseAuth: appAdmin,
	}
	authUserService := &services.AuthService{
		Repo: userRepo,
		FirebaseAuth: appUser,
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
	}

	// Init GIN
	r := gin.Default()

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
	routes.SetupRoutes(r, authAdminController, authUserController, userController)

	// Run server
	r.Run(":8080")
}
