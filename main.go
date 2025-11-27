package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/config"
	"github.com/muhammadfarrasfajri/login-google/database"
	routes "github.com/muhammadfarrasfajri/login-google/routers"
)

func main() {
	database.Connect()
	config.InitFirebase() 

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	routes.SetupRoutes(r)

	r.Run("0.0.0.0:8080")
}
