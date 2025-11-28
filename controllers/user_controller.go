package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/login-google/services"
)

type UserController struct {
	UserService *services.UserService
}

// GET /users
func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.UserService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get users",
		"users":   users,
	})
}

// GET /users/:id
func (c *UserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.UserService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get user",
		"user":    user,
	})
}

// PUT /users/:id
func (c *UserController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var body struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Role      string `json:"role"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := c.UserService.Update(id, body.Name, body.Email, body.AvatarURL, body.Role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update success",
		"user":    user,
	})
}

// DELETE /users/:id
func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.UserService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete success",
	})
}
