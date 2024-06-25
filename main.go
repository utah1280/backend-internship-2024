package main

import (
	"github.com/utah1280/backend-internship-2024/database/postgres"
	"github.com/utah1280/backend-internship-2024/internal/handlers/category"
	"github.com/utah1280/backend-internship-2024/internal/handlers/contact"
	"github.com/utah1280/backend-internship-2024/internal/server"
	"github.com/utah1280/backend-internship-2024/internal/storage"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			postgres.NewPostgresConnection,
			storage.NewCategoryStorage,
			storage.NewContactStorage,
			category.NewCategoryHandler,
			contact.NewContactHandler,
		),
		fx.Invoke(server.NewFiberServer),
	).Run()
}
