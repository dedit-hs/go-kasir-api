package services

import (
	"go-kasir-api/models"
	"go-kasir-api/repositories"
)

type CategoryServices struct {
	repository *repositories.CategoryRepositories
}

func NewCategoryServices(repository *repositories.CategoryRepositories) *CategoryServices {
	return &CategoryServices{repository: repository}
}

func (s *CategoryServices) GetAll() ([]models.Category, error) {
	return s.repository.GetAll()
}

func (s *CategoryServices) Create(category models.Category) (models.Category, error) {
	return s.repository.Create(category)
}

func (s *CategoryServices) GetByID(id int) (models.Category, error) {
	return s.repository.GetByID(id)
}

func (s *CategoryServices) Update(id int, category models.Category) (models.Category, error) {
	return s.repository.Update(id, category)
}

func (s *CategoryServices) Delete(id int) error {
	return s.repository.Delete(id)
}
