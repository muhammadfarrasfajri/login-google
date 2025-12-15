package repository

import (
	"database/sql"

	"github.com/muhammadfarrasfajri/login-google/models"
)

type PaymentRepository interface {
	// Update paramater: Terima Transaction header DAN slice of Details
	Save(transaction *models.Transaction, details []models.TransactionDetail) error
	UpdateStatus(orderID string, status string) error
	FindByOrderID(orderID string) (*models.Transaction, error)
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Save(t *models.Transaction, details []models.TransactionDetail) error {
	// 1. Mulai Database Transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// 2. Insert ke Tabel Transaction (Header)
	queryTx := `
		INSERT INTO transactions (order_id, user_id, amount, status, payment_type, payment_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(queryTx, t.OrderID, t.UserID, t.Amount, t.Status, t.PaymentType, t.PaymentUrl, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		tx.Rollback() // Batalkan jika gagal
		return err
	}

	// 3. Loop & Insert ke Tabel Transaction Details (Items)
	queryDetail := `
		INSERT INTO transaction_details (order_id, product_name, price, quantity, sub_total)
		VALUES (?, ?, ?, ?, ?)
	`
	for _, item := range details {
		_, err = tx.Exec(queryDetail, item.OrderID, item.ProductName, item.Price, item.Quantity, item.SubTotal)
		if err != nil {
			tx.Rollback() // Batalkan semua jika satu item gagal
			return err
		}
	}

	// 4. Commit (Simpan Permanen)
	return tx.Commit()
}

func (r *paymentRepository) UpdateStatus(orderID string, status string) error {
	query := `UPDATE transactions SET status = ? WHERE order_id = ?`
	_, err := r.db.Exec(query, status, orderID)
	return err
}

// Update FindByOrderID agar mengambil details juga (Optional tapi bagus)
func (r *paymentRepository) FindByOrderID(orderID string) (*models.Transaction, error) {
	// Query Header
	query := `SELECT id, order_id, user_id, amount, status, payment_type, payment_url, created_at, updated_at FROM transactions WHERE order_id = ?`
	row := r.db.QueryRow(query, orderID)

	var t models.Transaction
	err := row.Scan(&t.ID, &t.OrderID, &t.UserID, &t.Amount, &t.Status, &t.PaymentType, &t.PaymentUrl, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Query Details
	rows, err := r.db.Query("SELECT id, product_name, price, quantity, sub_total FROM transaction_details WHERE order_id = ?", orderID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var d models.TransactionDetail
			rows.Scan(&d.ID, &d.ProductName, &d.Price, &d.Quantity, &d.SubTotal)
			d.OrderID = orderID
			t.Details = append(t.Details, d) // Masukkan ke struct parent
		}
	}

	return &t, nil
}
