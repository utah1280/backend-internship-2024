package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Category struct {
	Id        int       `json:"id" db:"id"`
	Label     string    `json:"label" db:"label"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type NewCategoryInput struct {
	Label string
}

type CategoryStorage struct {
	DB *sqlx.DB
}

func NewCategoryStorage(DB *sqlx.DB) *CategoryStorage {
	return &CategoryStorage{DB: DB}
}

func GetCategoryIdByLabel(DB *sqlx.DB, label string) (int, error) {
	var id int

	stmt := "SELECT id FROM categories WHERE label = $1"
	err := DB.QueryRow(stmt, label).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("category '%s' does not exist", label)
		}
		return 0, fmt.Errorf("error getting category ID: %v", err)
	}
	return id, nil
}

func (storage *CategoryStorage) AddCategory(data NewCategoryInput) (int, error) {
	var id int

	checkStmt := "SELECT id FROM categories WHERE label = $1"
	err := storage.DB.QueryRow(checkStmt, data.Label).Scan(&id)

	if err == nil {
		return 0, fmt.Errorf("category '%s' already exists", data.Label)
	}

	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("error checking category existence: %v", err)
	}

	insertStmt := "INSERT INTO categories (label) VALUES ($1) RETURNING id"
	err = storage.DB.QueryRow(insertStmt, data.Label).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error adding category: %v", err)
	}

	return id, nil
}

func (storage *CategoryStorage) GetCategoryList() ([]Category, error) {
	var list []Category

	stmt := "SELECT * FROM categories"
	err := storage.DB.Select(&list, stmt)
	if err != nil {
		return nil, fmt.Errorf("error fetching category list: %v", err)
	}

	return list, nil
}

func (storage *CategoryStorage) DeleteCategory(id int) error {
	deleteStmt := "DELETE FROM categories WHERE id = $1"
	if _, err := storage.DB.Exec(deleteStmt, id); err != nil {
		return fmt.Errorf("error deleting category: %v", err)
	}
	return nil
}
