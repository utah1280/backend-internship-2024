package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Contact struct {
	Id         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Phone      string    `json:"phone" db:"phone"`
	Email      string    `json:"email" db:"email"`
	Address    string    `json:"address" db:"address"`
	CategoryId int       `json:"category_id" db:"category_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type NewContactInput struct {
	Name    string
	Phone   string
	Email   string
	Address string
	Label   string
}

type ContactStorage struct {
	DB *sqlx.DB
}

func NewContactStorage(DB *sqlx.DB) *ContactStorage {
	return &ContactStorage{DB: DB}
}

func (storage *ContactStorage) CreateContact(data NewContactInput) (int, error) {
	categoryId, err := GetCategoryIdByLabel(storage.DB, data.Label)
	if err != nil {
		return 0, fmt.Errorf("error fetching category id: %v", err)
	}

	var id int
	checkStmt := "SELECT id FROM contacts WHERE email = $1"
	err = storage.DB.QueryRow(checkStmt, data.Email).Scan(&id)

	switch {
	case err == sql.ErrNoRows:
		insertStmt := `
			INSERT INTO contacts (name, phone, email, address, category_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`
		err = storage.DB.QueryRow(insertStmt, data.Name, data.Phone, data.Email, data.Address, categoryId).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("error creating contact: %v", err)
		}
		return id, nil
	case err != nil:
		return 0, fmt.Errorf("error checking email existence: %v", err)
	default:
		return 0, fmt.Errorf("email '%s' already exists", data.Email)
	}
}

func (storage *ContactStorage) GetContact(id int) (Contact, error) {
	var contact Contact
	selectStmt := "SELECT id, name, phone, email, address, category_id, created_at FROM contacts WHERE id = $1"
	if err := storage.DB.Get(&contact, selectStmt, id); err != nil {
		return contact, fmt.Errorf("error fetching contact: %v", err)
	}
	return contact, nil
}

func (storage *ContactStorage) DeleteContact(id int) error {
	deleteStmt := "DELETE FROM contacts WHERE id = $1"
	if _, err := storage.DB.Exec(deleteStmt, id); err != nil {
		return fmt.Errorf("error deleting contact: %v", err)
	}
	return nil
}
