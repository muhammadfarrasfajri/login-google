package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/config"
	"github.com/muhammadfarrasfajri/login-google/database"
	"github.com/muhammadfarrasfajri/login-google/repository"
	service "github.com/muhammadfarrasfajri/login-google/services"
)

// GoogleAuth handles login with Firebase ID token
func GoogleAuth(c *gin.Context) {
	// Ambil token dari header, hapus "Bearer " jika ada
	idToken := c.GetHeader("Authorization")
	if idToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}
	if len(idToken) > 7 && idToken[:7] == "Bearer " {
		idToken = idToken[7:]
	}

	// Ambil device info dari body JSON
	var req struct {
		DeviceInfo string `json:"device_info"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userRepo := &repository.UserRepository{DB: database.DB}
	authService := service.NewAuthService(userRepo, config.FirebaseAuth)

	// Gunakan Remote IP dari request
	ipAddress := c.ClientIP()

	loginResp, err := authService.Login(idToken, req.DeviceInfo, ipAddress)
	if err != nil {
		if err == service.ErrUserNotRegistered {
			c.JSON(http.StatusForbidden, gin.H{
				"error":       err.Error(),
				"canRegister": true,
			})
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login success",
		"user":         loginResp.User,
		"sessionToken": loginResp.SessionToken,
		"expiresAt":    loginResp.ExpiresAt,
	})
}

// RegisterGoogle handles user registration from Firebase
func RegisterGoogle(c *gin.Context) {
	// Ambil token dari header, hapus "Bearer " jika ada
	idToken := c.GetHeader("Authorization")
	if idToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}
	if len(idToken) > 7 && idToken[:7] == "Bearer " {
		idToken = idToken[7:]
	}

	var req struct {
		DeviceInfo string `json:"device_info"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userRepo := &repository.UserRepository{DB: database.DB}
	authService := service.NewAuthService(userRepo, config.FirebaseAuth)

	ipAddress := c.ClientIP()

	loginResp, err := authService.Register(idToken, req.DeviceInfo, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Register success",
		"user":         loginResp.User,
		"sessionToken": loginResp.SessionToken,
		"expiresAt":    loginResp.ExpiresAt,
	})
}
