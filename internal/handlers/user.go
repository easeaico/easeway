package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"math/rand"
	"net/http"
	"net/mail"
	"time"

	_ "embed"

	"github.com/easeaico/easeway/internal/config"
	"github.com/easeaico/easeway/internal/store"
	"github.com/easeaico/easeway/internal/views/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

//go:embed verification_code.html
var verificationCodeTemplate string

func (u UserHandler) LoginPage(c echo.Context) error {
	sessionCookie, err := c.Cookie("session")
	if errors.Is(err, http.ErrNoCookie) {
		return user.LoginPage().Render(c.Request().Context(), c.Response().Writer)
	}

	if err != nil {
		slog.Error("session cookie not found", slog.Any("error", err))
		return err
	}

	sessionID := sessionCookie.Value
	_, err = u.queries.GetUserBySessionID(c.Request().Context(), pgtype.Text{
		String: sessionID,
		Valid:  true,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return user.LoginPage().Render(c.Request().Context(), c.Response().Writer)
	}
	if err != nil {
		slog.Error("get user by session error", slog.Any("error", err))
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (u UserHandler) DoLogin(c echo.Context) error {
	email := c.FormValue("signin-email")
	_, err := mail.ParseAddress(email)
	if err != nil {
		slog.Error("parse email address error", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "邮箱格式不正确")
	}

	ctx := c.Request().Context()

	user, err := u.queries.GetUserByEmail(ctx, email)
	if err != nil {
		slog.Error("query user by email error", slog.Any("error", err))
		return err
	}

	if time.Since(user.VerificationAt) > 5*time.Minute {
		return echo.NewHTTPError(http.StatusForbidden, "验证码已过期")
	}

	verificationCode := c.FormValue("signin-verifycode")
	if user.VerificationCode != verificationCode {
		return echo.NewHTTPError(http.StatusForbidden, "验证码不正确")
	}

	sessionID := generateSessionID()
	err = u.queries.UpdateSessionID(ctx, store.UpdateSessionIDParams{
		ID: user.ID,
		SessionID: pgtype.Text{
			String: sessionID,
			Valid:  true,
		},
	})
	if err != nil {
		slog.Error("update user session error", slog.Any("error", err))
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = sessionID
	cookie.Expires = time.Now().Add(10 * time.Hour)
	c.SetCookie(cookie)

	c.Response().Header().Set("HX-Redirect", "/")
	return nil
}

func (u UserHandler) SendVerification(c echo.Context) error {
	email := c.FormValue("signin-email")
	_, err := mail.ParseAddress(email)
	if err != nil {
		slog.Error("parse email address error", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "邮箱格式不正确")
	}

	t, err := template.New("verification_code").Parse(verificationCodeTemplate)
	if err != nil {
		slog.Error("parse message template error", slog.Any("error", err))
		return err
	}

	ctx := c.Request().Context()

	verificationCode := generateCode()
	user, err := u.queries.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		_, err := u.queries.CreateUser(ctx, store.CreateUserParams{
			Email:            email,
			VerificationCode: verificationCode,
		})
		if err != nil {
			slog.Error("create user by email error", slog.Any("error", err))
			return err
		}
	} else if err != nil {
		slog.Error("query user by email error", slog.Any("error", err))
		return err
	} else {
		err := u.queries.UpdateVerificationCode(ctx, store.UpdateVerificationCodeParams{
			ID:               user.ID,
			VerificationCode: verificationCode,
		})
		if err != nil {
			slog.Error("update user verification code error", slog.Any("error", err))
			return err
		}
	}

	data := struct {
		Code string
	}{
		Code: verificationCode,
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		slog.Error("render message template error", slog.Any("error", err))
		return err
	}

	content := buf.String()
	slog.Info("email template render", slog.String("content", content))

	m := gomail.NewMessage()
	cfg := u.config.Email
	m.SetHeader("From", cfg.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", cfg.Subject)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(cfg.ProviderHost, cfg.ProviderPort, cfg.From, cfg.Pass)
	if err := d.DialAndSend(m); err != nil {
		slog.Error("send email message error", slog.Any("error", err))
		return err
	}

	return nil
}

func generateCode() string {
	code := fmt.Sprintf("%06d", rand.Intn(10000))
	return code
}

func generateSessionID() string {
	return uuid.New().String()
}
