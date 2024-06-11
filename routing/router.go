package routing

import (
	"github.com/gofiber/fiber/v2"
	"go-pioneer/routing/app"
)

func Setup(f *fiber.App) {
	api := f.Group("/api")

	api.Post("/test", app.HandleTest)
	api.Get("/uploadPre", app.GetPreSignedUrl)
	api.Get("/downloadPre", app.DownloadPreSignedUrl)
}
