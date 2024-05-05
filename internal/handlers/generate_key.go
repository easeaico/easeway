package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/easeaico/easeway/internal/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type APIKeyHandler struct {
	queries *store.Queries
}

func NewAPIKeyHandler(queries *store.Queries) *APIKeyHandler {
	return &APIKeyHandler{
		queries: queries,
	}
}

func (k APIKeyHandler) GenerateNewKey(c echo.Context) error {
	key := uuid.New().String()
	key = strings.ReplaceAll(key, "-", "")
	key = fmt.Sprintf("ea-%s", key)
	_, err := k.queries.GetAPIKey(c.Request().Context(), key)
	if errors.Is(err, pgx.ErrNoRows) {
		_, err = k.queries.CreateAPIKey(c.Request().Context(), store.CreateAPIKeyParams{
			UserID: 1,
			Name:   "first key",
			Key:    key,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return nil
}
