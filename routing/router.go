package routing

import (
	"AlienegraGeek/go-pioneer/routing/app"
	"github.com/gofiber/fiber/v2"
)

func Setup(f *fiber.App) {
	appApi := f.Group("/app")

	appApi.Post("/test", app.HandleTest)
}
