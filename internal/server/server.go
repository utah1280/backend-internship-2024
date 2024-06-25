package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	_ "github.com/utah1280/backend-internship-2024/docs"
	"github.com/utah1280/backend-internship-2024/internal/handlers/category"
	"github.com/utah1280/backend-internship-2024/internal/handlers/contact"
	"go.uber.org/fx"
)

func NewFiberServer(lc fx.Lifecycle, contactHandlers *contact.ContactHandler, categoryHandlers *category.CategoryHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 4,
		WriteTimeout: time.Second * 4,
		IdleTimeout:  time.Second * 60,
	})
	app.Use(logger.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	contactGroup := app.Group("/contacts")
	contactGroup.Post("/new-contact", contactHandlers.CreateContact)
	contactGroup.Get("/get-contact/:id", contactHandlers.GetContact)
	contactGroup.Delete("/delete-contact/:id", contactHandlers.DeleteContact)
	contactGroup.Get("/get-contacts", contactHandlers.GetContacts)

	categoryGroup := app.Group("/categories")
	categoryGroup.Post("/add-category", categoryHandlers.AddCategory)
	categoryGroup.Get("/get-categories", categoryHandlers.GetCategoryList)
	categoryGroup.Delete("/delete-category/:id", categoryHandlers.DeleteCategory)
	categoryGroup.Patch("/update-category/:id", categoryHandlers.UpdateCategoryLabel)
	categoryGroup.Get("/get-category/:id", categoryHandlers.GetCategory)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting fiber server on port 8080")
			go func() {
				if err := app.Listen(":8080"); err != nil {
					fmt.Printf("Error starting server: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	return app
}
