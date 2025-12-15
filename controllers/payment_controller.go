package controllers

import (
	"net/http"

	"github.com/muhammadfarrasfajri/login-google/models"
	"github.com/muhammadfarrasfajri/login-google/services"

	// Sesuaikan dengan framework yang dipakai, misal "github.com/gofiber/fiber/v2"
	// Di sini saya pakai pseudo-code standar Go HTTP handler / Gin style
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	service services.PaymentService
}

func NewPaymentController(service services.PaymentService) *PaymentController {
	return &PaymentController{service}
}

// CreatePayment menangani request pembuatan QRIS
func (c *PaymentController) CreatePayment(ctx *gin.Context) {
	var req models.PaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil User ID dari Token/Middleware (Contoh)
	userID := ctx.GetInt("user_id")
	// userID := strconv.Itoa(userID1)

	resp, err := c.service.CreateQrisTransaction(req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "QRIS generated successfully",
		"data":    resp,
	})
}

// MidtransCallback menangani Webhook dari Midtrans
func (c *PaymentController) MidtransCallback(ctx *gin.Context) {
	var notificationPayload map[string]interface{}

	if err := ctx.ShouldBindJSON(&notificationPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.HandleNotification(notificationPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process notification"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
