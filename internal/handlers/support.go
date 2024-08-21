package handlers

import (
	"github.com/easeaico/easeway/internal/store"
	"github.com/easeaico/easeway/internal/views/support"
	"github.com/labstack/echo/v4"
)

type SupportHandler struct{}

func NewSupportHandler() *SupportHandler {
	return &SupportHandler{}
}

func (h SupportHandler) HomePage(c echo.Context) error {
	ctx := c.Request().Context()
	writer := c.Response().Writer
	user, ok := c.Get("user").(*store.User)
	if ok {
		return support.Home(user.Email, true).Render(ctx, writer)
	} else {
		return support.Home("", false).Render(ctx, writer)
	}
}
