package repositories

import (
	"database/sql"
	"errors"
	"go-kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name, description FROM categories ORDER BY id ASC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoryRepository) Create(category models.Category) (models.Category, error) {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	row := r.db.QueryRow(query, category.Name, category.Description)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return models.Category{}, err
	}
	category.ID = id
	return category, nil
}

func (r *CategoryRepository) GetByID(id int) (models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"
	row := r.db.QueryRow(query, id)
	var category models.Category
	err := row.Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func (r *CategoryRepository) Update(id int, category models.Category) (models.Category, error) {
	query := "UPDATE categories SET name = $2, description = $3 WHERE id = $1 RETURNING id"
	result, err := r.db.Exec(query, id, category.Name, category.Description)
	if err != nil {
		return models.Category{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return models.Category{}, err
	}
	if rows == 0 {
		return models.Category{}, errors.New("category not found")
	}

	return category, nil
}

func (r *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("category not found")
	}
	return nil
}
