package handlers

import (
	"github.com/easeaico/easeway/internal/views/support"
	"github.com/labstack/echo/v4"
)

type SupportHandler struct{}

func NewSupportHandler() *SupportHandler {
	return &SupportHandler{}
}

func (h SupportHandler) HomePage(c echo.Context) error {
	return support.HomePage().Render(c.Request().Context(), c.Response().Writer)
}
