package handlers

import (
	"database/sql"
	"naiimage/backend/internal/nai"
	"naiimage/backend/internal/store"
)

// Handler 持有所有依赖，通过 fiber.Locals 或闭包注入到各路由。
type Handler struct {
	DB       *sql.DB
	Defaults store.Settings
	Client   *nai.Client
}
