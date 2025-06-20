package handlers

import (
	"github.com/easeaico/easeway/internal/store"
	"github.com/easeaico/easeway/internal/views/home"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
	queries *store.Queries
}

func NewHomeHandler(queries *store.Queries) *HomeHandler {
	return &HomeHandler{
		queries: queries,
	}
}

func (h HomeHandler) HomePage(c echo.Context) error {
	ctx := c.Request().Context()
	writer := c.Response().Writer
	return home.Home("", false).Render(ctx, writer)
}
