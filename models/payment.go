package models

import "time"

// PaymentRequest dari Frontend
type PaymentRequest struct {
	// Amount int64  `json:"amount"` <-- HAPUS INI (Sebaiknya hitung otomatis dari total barang)
	User_id int           `json:"user_id"`
	Email   string        `json:"email"`
	Name    string        `json:"name"`
	Items   []ItemRequest `json:"items"` // List barang yang dibeli
}

type ItemRequest struct {
	ProductName string `json:"product_name"`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
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
	UserID      int       `json:"user_id"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
	PaymentType string    `json:"payment_type"`
	PaymentUrl  string    `json:"payment_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Optional: Untuk keperluan query response nanti
	Details []TransactionDetail `json:"details,omitempty"`
}

// 3. Struct Database (Child)
type TransactionDetail struct {
	ID          int64  `json:"id"`
	OrderID     string `json:"order_id"`
	ProductName string `json:"product_name"`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
	SubTotal    int64  `json:"sub_total"`
}
