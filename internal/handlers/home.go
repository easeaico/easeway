package handlers

import (
	"github.com/easeaico/easeway/internal/views/home"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h HomeHandler) HomePage(c echo.Context) error {
	return home.HomePage("", false).Render(c.Request().Context(), c.Response().Writer)
}
