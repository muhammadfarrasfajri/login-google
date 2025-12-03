package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/services"
)

type AuthController struct {
	AuthService *services.AuthService
}

func (c *AuthController) RegisterUser(ctx *gin.Context) {
	var body struct {
		IDToken string `json:"id_token"`
		Name    string `json:"name"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := c.AuthService.Register(body.IDToken, body.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Register success",
		"user":    user,
	})
}

func (c *AuthController) LoginUser(ctx *gin.Context) {
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
