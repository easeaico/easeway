package handlers

import (
	"github.com/easeaico/easeway/internal/store"
	"github.com/easeaico/easeway/internal/views/member"
	"github.com/labstack/echo/v4"
)

type MemberHandler struct{}

func NewMemberHandler() *MemberHandler {
	return &MemberHandler{}
}

func (h MemberHandler) HomePage(c echo.Context) error {
	ctx := c.Request().Context()
	writer := c.Response().Writer
	user, ok := c.Get("user").(*store.User)
	if ok {
		return member.Home(user.Email, true).Render(ctx, writer)
	} else {
		return member.Home("", false).Render(ctx, writer)
	}
}
