package models

import "time"

// PaymentRequest dari Frontend
type PaymentRequest struct {
	Amount int64  `json:"amount"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// PaymentResponse ke Frontend
type PaymentResponse struct {
	OrderID    string `json:"order_id"`
	Amount     string `json:"amount"`
	QRImageUrl string `json:"qr_image_url"`
	Status     string `json:"status"`
}

// Transaction struct (Representasi Tabel)
type Transaction struct {
	ID          int64     `json:"id"`
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
	PaymentType string    `json:"payment_type"`
	PaymentUrl  string    `json:"payment_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
