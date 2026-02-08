package repositories

import (
	"database/sql"
	"errors"
	"go-kasir-api/models"
)

type ProductRepositories struct {
	db *sql.DB
}

func NewProductRepositories(db *sql.DB) *ProductRepositories {
	return &ProductRepositories{db: db}
}

func (r *ProductRepositories) GetAll(name string) ([]models.Product, error) {
	args := []interface{}{}
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id, 
	          c.id, c.name, c.description 
	          FROM products p 
	          JOIN categories c ON p.category_id = c.id`
	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}
	query += " ORDER BY p.id ASC"
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var category models.Category
		err := rows.Scan(
			&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID,
			&category.ID, &category.Name, &category.Description,
		)
		if err != nil {
			return nil, err
		}
		product.Category = &category
		products = append(products, product)
	}
	return products, nil
}

func (r *ProductRepositories) Create(product models.Product) (models.Product, error) {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	row := r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return models.Product{}, err
	}
	product.ID = id
	return product, nil
}

func (r *ProductRepositories) GetByID(id int) (models.Product, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id,
	          c.id, c.name, c.description
	          FROM products p
	          JOIN categories c ON p.category_id = c.id
	          WHERE p.id = $1`
	row := r.db.QueryRow(query, id)
	var product models.Product
	var category models.Category
	err := row.Scan(
		&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID,
		&category.ID, &category.Name, &category.Description,
	)
	if err != nil {
		return models.Product{}, err
	}
	product.Category = &category
	return product, nil
}

func (r *ProductRepositories) Update(id int, product models.Product) (models.Product, error) {
	query := "UPDATE products SET name = $2, price = $3, stock = $4, category_id = $5 WHERE id = $1 RETURNING id"
	result, err := r.db.Exec(query, id, product.Name, product.Price, product.Stock, product.CategoryID)
	if err != nil {
		return models.Product{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return models.Product{}, err
	}
	if rows == 0 {
		return models.Product{}, errors.New("product not found")
	}

	return product, nil
}

func (r *ProductRepositories) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}
