package services

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/muhammadfarrasfajri/login-google/models"
	"github.com/muhammadfarrasfajri/login-google/repository"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type PaymentService interface {
	CreateQrisTransaction(req models.PaymentRequest, userID string) (*models.PaymentResponse, error)
	HandleNotification(notificationPayload map[string]interface{}) error
}

type paymentService struct {
	repo   repository.PaymentRepository
	client coreapi.Client
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	// Setup Midtrans Client
	var client coreapi.Client
	env := midtrans.Sandbox
	if os.Getenv("MIDTRANS_IS_PRODUCTION") == "true" {
		env = midtrans.Production
	}
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), env)

	return &paymentService{repo: repo, client: client}
}

func (s *paymentService) CreateQrisTransaction(req models.PaymentRequest, userID string) (*models.PaymentResponse, error) {
	orderID := "ORDER-" + uuid.New().String()

	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeQris,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(req.Amount),
		},
		Qris: &coreapi.QrisDetails{Acquirer: "gopay"},
	}

	resp, err := s.client.ChargeTransaction(chargeReq)
	if err != nil {
		return nil, err
	}

	var qrImageUrl string
	for _, action := range resp.Actions {
		if action.Name == "generate-qr-code" {
			qrImageUrl = action.URL
		}
	}

	// Buat Model Transaction
	now := time.Now()
	tx := &models.Transaction{
		OrderID:     orderID,
		UserID:      userID,
		Amount:      int64(req.Amount),
		Status:      resp.TransactionStatus,
		PaymentType: "qris",
		PaymentUrl:  qrImageUrl,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Simpan via Repo (Tanpa GORM)
	if err := s.repo.Save(tx); err != nil {
		return nil, err
	}

	return &models.PaymentResponse{
		OrderID:    orderID,
		Amount:     strconv.FormatInt(int64(req.Amount), 10),
		QRImageUrl: qrImageUrl,
		Status:     resp.TransactionStatus,
	}, nil
}

func (s *paymentService) HandleNotification(payload map[string]interface{}) error {
	orderID, _ := payload["order_id"].(string)
	transactionStatus, _ := payload["transaction_status"].(string)
	fraudStatus, _ := payload["fraud_status"].(string)
	statusCode, _ := payload["status_code"].(string)
	grossAmount, _ := payload["gross_amount"].(string)
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")

	// 1. Ambil Signature Key dari Payload Midtrans
	signatureKey, _ := payload["signature_key"].(string)

	// 2. Generate Signature Sendiri untuk Validasi
	// Rumus Midtrans: SHA512(order_id + status_code + gross_amount + ServerKey)
	input := orderID + statusCode + grossAmount + serverKey
	hasher := sha512.New()
	hasher.Write([]byte(input))
	expectedSignature := hex.EncodeToString(hasher.Sum(nil))

	// 3. Bandingkan
	if signatureKey != expectedSignature {
		return errors.New("invalid signature key") // Tolak request palsu
	}

	var newStatus string = "pending"

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			newStatus = "challenge"
		} else if fraudStatus == "accept" {
			newStatus = "success"
		}
	} else if transactionStatus == "settlement" {
		newStatus = "success"
	} else if transactionStatus == "deny" || transactionStatus == "cancel" || transactionStatus == "expire" {
		newStatus = "failed"
	}

	// Update Status via Repo (Query SQL Update)
	return s.repo.UpdateStatus(orderID, newStatus)
}
