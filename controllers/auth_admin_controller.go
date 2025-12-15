package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *AuthController) RegisterAdmin(ctx *gin.Context) {
	var body struct {
		IDToken string `json:"id_token"`
		Name    string `json:"name"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	admin, err := c.AuthService.Register(body.IDToken, body.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Register success",
		"user":    admin,
	})
}

func (c *AuthController) LoginAdmin(ctx *gin.Context) {
	var req struct {
		IDToken    string `json:"id_token"`
		DeviceInfo string `json:"device_info"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ip := ctx.ClientIP()

	result, err := c.AuthService.Login(req.IDToken, req.DeviceInfo, ip)

	if err != nil {
		ctx.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *AuthController) RefreshTokenAdmin(ctx *gin.Context) {

	refreshToken, err := ctx.Cookie("refresh_token")
	fmt.Println(refreshToken)
	log.Println("Error retrieving refresh token from cookie:", refreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}
	result, err := c.AuthService.RefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *AuthController) LogoutAdmin(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	err := c.AuthService.Logout(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logout success"})
}
