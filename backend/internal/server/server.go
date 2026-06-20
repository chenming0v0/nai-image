package server

import (
	"errors"
	"log"
	"time"

	"naiimage/backend/internal/config"
	"naiimage/backend/internal/db"
	"naiimage/backend/internal/handlers"
	"naiimage/backend/internal/nai"
	"naiimage/backend/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Run(cfg config.Config) error {
	if err := db.Init(cfg.DBPath); err != nil {
		return err
	}

	defaults := store.Settings{
		UpstreamBaseURL: cfg.UpstreamBaseURL,
		UpstreamAPIKey:  cfg.UpstreamAPIKey,
		DefaultModel:    cfg.DefaultModel,
		RequestTimeout:  int(cfg.RequestTimeout / time.Second),
	}
	if defaults.RequestTimeout <= 0 {
		defaults.RequestTimeout = 180
	}

	timeout := cfg.RequestTimeout
	if timeout <= 0 {
		timeout = 180 * time.Second
	}

	h := &handlers.Handler{
		DB:       db.DB,
		Defaults: defaults,
		Client:   nai.NewClient(timeout),
	}

	app := fiber.New(fiber.Config{
		BodyLimit: int(cfg.MaxImageBytes),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var fe *fiber.Error
			if errors.As(err, &fe) {
				code = fe.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api := app.Group("/api")
	api.Get("/health", handlers.HealthHandler)
	api.Get("/settings", h.GetSettingsHandler)
	api.Put("/settings", h.UpdateSettingsHandler)
	api.Get("/models", h.GetModelsHandler)
	api.Post("/generate", h.GenerateHandler)
	api.Get("/tasks", h.ListTasksHandler)
	api.Delete("/tasks", h.DeleteAllTasksHandler)
	api.Get("/tasks/:id", h.GetTaskHandler)
	api.Delete("/tasks/:id", h.DeleteTaskHandler)
	api.Get("/images/:id/meta", h.GetImageMetaHandler)
	api.Get("/images/:id", h.GetImageHandler)

	addr := ":" + cfg.Port
	log.Printf("nai-image backend listening on http://127.0.0.1%s", addr)
	return app.Listen(addr)
}
