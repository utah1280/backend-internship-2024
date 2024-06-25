package category

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/utah1280/backend-internship-2024/internal/storage"
)

type CategoryHandler struct {
	Storage *storage.CategoryStorage
}

func NewCategoryHandler(storage *storage.CategoryStorage) *CategoryHandler {
	return &CategoryHandler{Storage: storage}
}

type basicResponse struct {
	Success bool `json:"success"`
}

type categoryRequest struct {
	Label string `json:"label"`
}

type categoryResponse struct {
	Id int `json:"id"`
}

// AddCategory swagger
// @Summary Create a new category
// @Description Create a new category with the given label
// @Tags Categories
// @Accept json
// @Produce json
// @Param body body categoryRequest true "Category details"
// @Success 200 {object} categoryResponse
// @Failure 400 {string} string "Bad Request"
// @Router /categories/add-category [post]
func (handler *CategoryHandler) AddCategory(ctx *fiber.Ctx) error {
	var body categoryRequest

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	id, err := handler.Storage.AddCategory(storage.NewCategoryInput{
		Label: body.Label,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	resp := categoryResponse{Id: id}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

type categoryListResponse struct {
	Categories []storage.Category `json:"categories"`
}

// GetCategoryList swagger
// @Summary Get list of categories
// @Description Retrieve a list of all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} categoryListResponse
// @Router /categories/get-categories [get]
func (handler *CategoryHandler) GetCategoryList(ctx *fiber.Ctx) error {
	categories, err := handler.Storage.GetCategoryList()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	resp := categoryListResponse{
		Categories: categories,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// DeleteCategory swagger
// @Summary Delete category
// @Description Delete category with the given id
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} basicResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /categories/delete-category/{id} [delete]
func (handler *CategoryHandler) DeleteCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	categoryId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid category ID")
	}

	err = handler.Storage.DeleteCategory(categoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).SendString("Category not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	res := basicResponse{
		Success: true,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

type updateCategoryLabelRequest struct {
	Label string `json:"label"`
}

// UpdateCategoryLabel swagger
// @Summary Update category label
// @Description Update the label of a category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param body body updateCategoryLabelRequest true "Category details"
// @Success 200 {object} basicResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /categories/update-category/{id} [patch]
func (handler *CategoryHandler) UpdateCategoryLabel(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req updateCategoryLabelRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	categoryID, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid category ID")
	}

	err = handler.Storage.UpdateCategoryLabel(categoryID, req.Label)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to update category label: %v", err))
	}

	resp := basicResponse{Success: true}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

type fetchCategoryRespones struct {
	Category storage.Category `json:"category"`
}

// GetCategory swagger
// @Summary Get a category by ID
// @Description Retrieve details of a category based on the provided ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} fetchCategoryRespones
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /categories/get-category/{id} [get]
func (handler *CategoryHandler) GetCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	categoryId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid category ID")
	}

	category, err := handler.Storage.GetCategory(categoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Category not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	resp := fetchCategoryRespones{
		Category: category,
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
