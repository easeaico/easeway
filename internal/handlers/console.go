package handlers

import (
	"errors"
	"fmt"
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
	return console.HomePage(user.Email, true).Render(c.Request().Context(), c.Response().Writer)
}

func (h ConsoleHandler) GenerateNewKey(c echo.Context) error {
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

	return nil
}

func generateAPIKey() string {
	key := uuid.New().String()
	key = strings.ReplaceAll(key, "-", "")
	return fmt.Sprintf("ea-%s", key)
}
