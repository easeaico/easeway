package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"net/mail"

	_ "embed"

	"github.com/easeaico/easeway/internal/config"
	"github.com/easeaico/easeway/internal/store"
	"github.com/easeaico/easeway/internal/views/user"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

type UserHandler struct {
	config  *config.Config
	queries *store.Queries
}

func NewUserHandler(config *config.Config, queries *store.Queries) *UserHandler {
	return &UserHandler{
		config:  config,
		queries: queries,
	}
}

// @embed verification_code.html
var verificationCodeTemplate string

func (u UserHandler) LoginPage(c echo.Context) error {
	return user.LoginPage().Render(c.Request().Context(), c.Response().Writer)
}

func (u UserHandler) DoLogin(c echo.Context) error {
	return nil
}

func (u UserHandler) SendVerificationCode(c echo.Context) error {
	email := c.FormValue("email")
	_, err := mail.ParseAddress(email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "邮箱格式不正确")
	}

	t, err := template.New("verification_code").Parse(verificationCodeTemplate)
	if err != nil {
		return err
	}

	data := struct {
		Code string
	}{
		Code: "1",
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	cfg := u.config.Email
	m.SetHeader("From", cfg.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", cfg.Subject)
	m.SetBody("text/html", buf.String())

	d := gomail.NewDialer(cfg.ProviderHost, cfg.ProviderPort, cfg.From, cfg.Pass)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// 生成6位验证码
func generateCode() string {
	code := fmt.Sprintf("%06d", rand.Intn(10000))
	return code
}
