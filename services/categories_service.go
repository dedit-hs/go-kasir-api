package services

import (
	"go-kasir-api/models"
	"go-kasir-api/repositories"
)

type CategoryService struct {
	repository *repositories.CategoryRepository
}

func NewCategoryService(repository *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repository: repository}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repository.GetAll()
}

func (s *CategoryService) Create(category models.Category) (models.Category, error) {
	return s.repository.Create(category)
}

func (s *CategoryService) GetByID(id int) (models.Category, error) {
	return s.repository.GetByID(id)
}

func (s *CategoryService) Update(id int, category models.Category) (models.Category, error) {
	return s.repository.Update(id, category)
}

func (s *CategoryService) Delete(id int) error {
	return s.repository.Delete(id)
}
