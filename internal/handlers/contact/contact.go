package contact

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/utah1280/backend-internship-2024/internal/storage"
)

type ContactHandler struct {
	Storage *storage.ContactStorage
}

func NewContactHandler(storage *storage.ContactStorage) *ContactHandler {
	return &ContactHandler{Storage: storage}
}

type basicResponse struct {
	Success bool `json:"success"`
}

type createContactRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Label   string `json:"label"`
}

type createContactResponse struct {
	Id int `json:"id"`
}

// CreateContact swagger
// @Summary Create a new contact
// @Description Create a new contact with the given details
// @Tags Contacts
// @Accept json
// @Produce json
// @Param body body createContactRequest true "Contact details"
// @Success 200 {object} createContactResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /contacts/new-contact [post]
func (handler *ContactHandler) CreateContact(ctx *fiber.Ctx) error {
	var body createContactRequest

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	id, err := handler.Storage.CreateContact(storage.NewContactInput{
		Name:    body.Name,
		Phone:   body.Phone,
		Email:   body.Email,
		Address: body.Address,
		Label:   body.Label,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	resp := createContactResponse{Id: id}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

type fetchContactResponse struct {
	Contact storage.Contact_ `json:"contact"`
}

// GetContact swagger
// @Summary Get a contact by ID
// @Description Retrieve details of a contact based on the provided ID
// @Tags Contacts
// @Accept json
// @Produce json
// @Param id path int true "Contact ID"
// @Success 200 {object} fetchContactResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /contacts/get-contact/{id} [get]
func (handler *ContactHandler) GetContact(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	contactId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid contact ID")
	}

	contact, err := handler.Storage.GetContact(contactId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).SendString("Contact not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	resp := fetchContactResponse{
		Contact: contact,
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// DeleteContact swagger
// @Summary Delete contact
// @Description Delete contact with the given id
// @Tags Contacts
// @Accept json
// @Produce json
// @Param id path int true "Contact ID"
// @Success 200 {object} basicResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /contacts/delete-contact/{id} [delete]
func (handler *ContactHandler) DeleteContact(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	contactId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid contact ID")
	}

	err = handler.Storage.DeleteContact(contactId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).SendString("Contact not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	res := basicResponse{
		Success: true,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

type contactListResponse struct {
	Contacts []storage.Contact_ `json:"contact"`
}

// GetContacts swagger
// @Summary Get list of contacts
// @Description Retrieve a list of contacts with optional filtering, sorting, and pagination
// @Tags Contacts
// @Accept json
// @Produce json
// @Param limit query int false "Limit results per page"
// @Param offset query int false "Offset results for pagination"
// @Param name query string false "Filter by contact name"
// @Param email query string false "Filter by contact email"
// @Param category query string false "Filter by category label"
// @Param sortDir query string false "Sort direction (ASC default)"
// @Success 200 {array} storage.Contact_
// @Router /contacts/get-contacts [get]
func (handler *ContactHandler) GetContacts(ctx *fiber.Ctx) error {
	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(ctx.Query("offset", "0"))
	if err != nil {
		offset = 0
	}

	name := ctx.Query("name", "")
	email := ctx.Query("email", "")
	category := ctx.Query("category", "")
	sortDir := ctx.Query("sortDir", "ASC")

	contacts, err := handler.Storage.GetContacts(limit, offset, name, email, category, sortDir)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	resp := contactListResponse{
		Contacts: contacts,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
