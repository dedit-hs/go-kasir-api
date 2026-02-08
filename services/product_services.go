package services

import (
	"go-kasir-api/models"
	"go-kasir-api/repositories"
)

type ProductServices struct {
	productRepo *repositories.ProductRepositories
}

func NewProductServices(productRepo *repositories.ProductRepositories) *ProductServices {
	return &ProductServices{productRepo: productRepo}
}

func (s *ProductServices) GetAll(name string) ([]models.Product, error) {
	return s.productRepo.GetAll(name)
}

func (s *ProductServices) Create(product models.Product) (models.Product, error) {
	return s.productRepo.Create(product)
}

func (s *ProductServices) GetByID(id int) (models.Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *ProductServices) Update(id int, product models.Product) (models.Product, error) {
	return s.productRepo.Update(id, product)
}

func (s *ProductServices) Delete(id int) error {
	return s.productRepo.Delete(id)
}
