package repository

import (
	"database/sql"
	"errors"

	"github.com/muhammadfarrasfajri/login-google/models"
)

type PaymentRepository interface {
	Save(transaction *models.Transaction) error
	UpdateStatus(orderID string, status string) error
	FindByOrderID(orderID string) (*models.Transaction, error)
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Save(t *models.Transaction) error {
	query := `
		INSERT INTO transactions (order_id, user_id, amount, status, payment_type, payment_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	// Gunakan $1, $2 dst jika menggunakan PostgreSQL, gunakan ? untuk MySQL
	_, err := r.db.Exec(query, t.OrderID, t.UserID, t.Amount, t.Status, t.PaymentType, t.PaymentUrl, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *paymentRepository) UpdateStatus(orderID string, status string) error {
	query := `UPDATE transactions SET status = ? WHERE order_id = ?`
	_, err := r.db.Exec(query, status, orderID)
	return err
}

func (r *paymentRepository) FindByOrderID(orderID string) (*models.Transaction, error) {
	query := `SELECT id, order_id user_id, amount, status, payment_type, payment_url, created_at, updated_at FROM transactions WHERE order_id = ?`

	row := r.db.QueryRow(query, orderID)

	var t models.Transaction
	// Scan urutannya harus sama persis dengan SELECT di atas
	err := row.Scan(
		&t.ID,
		&t.OrderID,
		&t.UserID,
		&t.Amount,
		&t.Status,
		&t.PaymentType,
		&t.PaymentUrl,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	return &t, nil
}
