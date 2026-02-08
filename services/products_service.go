package services

import (
	"go-kasir-api/models"
	"go-kasir-api/repositories"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.productRepo.GetAll(name)
}

func (s *ProductService) Create(product models.Product) (models.Product, error) {
	return s.productRepo.Create(product)
}

func (s *ProductService) GetByID(id int) (models.Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *ProductService) Update(id int, product models.Product) (models.Product, error) {
	return s.productRepo.Update(id, product)
}

func (s *ProductService) Delete(id int) error {
	return s.productRepo.Delete(id)
}
