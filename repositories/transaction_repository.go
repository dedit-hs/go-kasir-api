package repositories

import (
	"database/sql"
	"errors"
	"go-kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productID, productPrice, stock int
		var productName string
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productID, &productName, &productPrice, &stock)
		if err != nil {
			return nil, errors.New("Product not found")
		}

		if stock < item.Quantity {
			return nil, errors.New("Insufficient stock")
		}

		// Calculate subtotal
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// Update stock
		_, err = tx.Exec("UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// Create transaction detail
		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// Insert details and capture IDs
	for i := range details {
		var detailID int
		err := tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id", transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		// Update the detail with IDs
		details[i].ID = detailID
		details[i].TransactionID = transactionID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
