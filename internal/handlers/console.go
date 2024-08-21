package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/easeaico/easeway/internal/store"
	"github.com/easeaico/easeway/internal/views/console"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type ConsoleHandler struct {
	queries *store.Queries
}

func NewConsoleHandler(queries *store.Queries) *ConsoleHandler {
	return &ConsoleHandler{
		queries: queries,
	}
}

func (h ConsoleHandler) HomePage(c echo.Context) error {
	user := c.Get("user").(*store.User)
	keys, err := h.queries.ListAPIKeys(c.Request().Context(), user.ID)
	if err != nil && err != pgx.ErrNoRows {
		slog.Error("list api keys error", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return console.HomePage(user.Email, true, keys).Render(c.Request().Context(), c.Response().Writer)
}

func (h ConsoleHandler) CreateKeyPage(c echo.Context) error {
	user := c.Get("user").(*store.User)
	return console.CreateKeyPage(user.Email, true).Render(c.Request().Context(), c.Response().Writer)
}

func (h ConsoleHandler) GenerateKey(c echo.Context) error {
	user := c.Get("user").(*store.User)

	name := c.FormValue("key_name")
	key := generateAPIKey()
	_, err := h.queries.GetAPIKey(c.Request().Context(), key)
	if errors.Is(err, pgx.ErrNoRows) {
		_, err = h.queries.CreateAPIKey(c.Request().Context(), store.CreateAPIKeyParams{
			UserID: user.ID,
			Name:   name,
			Key:    key,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/console/home")
}

func generateAPIKey() string {
	key := uuid.New().String()
	key = strings.ReplaceAll(key, "-", "")
	return fmt.Sprintf("ea-%s", key)
}
