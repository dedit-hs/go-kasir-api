package services

import (
	"go-kasir-api/models"
	"go-kasir-api/repositories"
)

type TransactionService struct {
	productRepo     *repositories.ProductRepository
	transactionRepo *repositories.TransactionRepository
}

func NewTransactionService(productRepo *repositories.ProductRepository, transactionRepo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{productRepo: productRepo, transactionRepo: transactionRepo}
}

func (s *TransactionService) Create(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.transactionRepo.Create(items)
}
