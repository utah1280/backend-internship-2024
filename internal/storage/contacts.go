package storage

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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

type Contact_ struct {
	Id         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Phone      string    `json:"phone" db:"phone"`
	Email      string    `json:"email" db:"email"`
	Address    string    `json:"address" db:"address"`
	CategoryId int       `json:"category_id" db:"category_id"`
	Category   string    `json:"category" db:"category"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
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

func (storage *ContactStorage) GetContact(id int) (Contact_, error) {
	var contact Contact_
	selectStmt := `
		SELECT c.id, c.name, c.phone, c.email, c.address, c.category_id, cat.label as category, c.created_at 
		FROM contacts c
		LEFT JOIN categories cat ON c.category_id = cat.id
		WHERE c.id = $1
	`
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

func (storage *ContactStorage) GetContacts(limit, offset int, name, email, category, sortDir string) ([]Contact_, error) {
	var contacts []Contact_
	stmt := `
		SELECT c.id, c.name, c.phone, c.email, c.address, c.category_id, cat.label as category, c.created_at 
		FROM contacts c
		LEFT JOIN categories cat ON c.category_id = cat.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if name != "" {
		stmt += " AND c.name ILIKE $" + strconv.Itoa(argCount)
		args = append(args, "%"+name+"%")
		argCount++
	}

	if email != "" {
		stmt += " AND c.email ILIKE $" + strconv.Itoa(argCount)
		args = append(args, "%"+email+"%")
		argCount++
	}

	if category != "" {
		stmt += `
			AND c.category_id IN (
				SELECT id FROM categories WHERE label ILIKE $` + strconv.Itoa(argCount) + `
			)
		`
		args = append(args, "%"+category+"%")
		argCount++
	}

	if sortDir != "" {
		stmt += " ORDER BY c.created_at " + sortDir
	}

	if limit > 0 {
		stmt += " LIMIT $" + strconv.Itoa(argCount)
		args = append(args, limit)
		argCount++
	}

	if offset > 0 {
		if limit <= 0 {
			return nil, fmt.Errorf("offset specified without limit")
		}
		stmt += " OFFSET $" + strconv.Itoa(argCount)
		args = append(args, offset)
		argCount++
	}

	if offset > 0 {
		if limit <= 0 {
			return nil, fmt.Errorf("offset specified without limit")
		}
		stmt += " OFFSET $" + strconv.Itoa(argCount)
		args = append(args, offset)
		argCount++
	}

	err := storage.DB.Select(&contacts, stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("error retrieving contacts: %v", err)
	}

	return contacts, nil
}

func (storage *ContactStorage) UpdateContact(id int, name, phone, email, address, category string) error {
	var temp int
	checkStmt := "SELECT id FROM contacts WHERE email = $1"
	err := storage.DB.QueryRow(checkStmt, email).Scan(&temp)

	if err == nil {
		return fmt.Errorf("email '%s' already exists", email)
	}

	stmt := "UPDATE contacts SET"
	args := []interface{}{}
	args = append(args, id)

	if name != "" {
		stmt += " name = $2,"
		args = append(args, name)
	}
	if phone != "" {
		stmt += " phone = $3,"
		args = append(args, phone)
	}
	if address != "" {
		stmt += " address = $4,"
		args = append(args, address)
	}
	if category != "" {
		categoryId, err := GetCategoryIdByLabel(storage.DB, category)
		if err != nil {
			return fmt.Errorf("error fetching category id: %v", err)
		}

		stmt += " category_id = $" + strconv.Itoa(len(args)+1) + ","
		args = append(args, categoryId)
	}
	if email != "" {
		stmt += " email = $" + strconv.Itoa(len(args)+1) + ","
		args = append(args, email)
	}

	stmt = strings.TrimSuffix(stmt, ",")
	stmt += " WHERE id = $1"

	_, err = storage.DB.Exec(stmt, args...)
	if err != nil {
		return fmt.Errorf("error updating contact: %v", err)
	}

	return nil
}
