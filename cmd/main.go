package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/easeaico/easeway/internal/config"
	"github.com/easeaico/easeway/internal/handlers"
	"github.com/easeaico/easeway/internal/spi"
	"github.com/easeaico/easeway/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	confFile := flag.String("f", "config.yaml", "config file path")
	flag.Parse()

	conf := config.NewConfig(*confFile)
	ctx := context.Background()
	spis := spi.NewSPIRegistry(ctx, conf)

	pool := store.NewDBTX(ctx, conf)
	defer pool.Close()
	queries := store.New(pool)

	e := echo.New()
	e.File("/favicon.ico", "assets/images/favicon.ico")
	e.Static("/assets", "assets")

	e.Use(middleware.Logger(), middleware.Recover(), middleware.CORS())

	setupRoutes(e, conf, spis, queries)

	if err := e.Start(fmt.Sprintf("%s:%d", conf.Server.IP, conf.Server.Port)); err != nil {
		e.Logger.Fatal(err)
	}
}

func setupRoutes(e *echo.Echo, conf *config.Config, spis *spi.SPIRegistry, queries *store.Queries) {
	homeHandler := handlers.NewHomeHandler(queries)
	openaiApiHandler := handlers.NewAPISvcHandler(spis, queries)
	anthropicApiHandler := handlers.NewAthropicApiHandler(conf, queries)
	userHandler := handlers.NewUserHandler(conf, queries)
	consoleHandler := handlers.NewConsoleHandler(queries)
	supportHandler := handlers.NewSupportHandler()
	memberHandler := handlers.NewMemberHandler()

	e.GET("/", homeHandler.HomePage)
	e.GET("/support", supportHandler.HomePage)
	e.GET("/member", memberHandler.HomePage)

	v1 := e.Group("/v1", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Validator: func(key string, c echo.Context) (bool, error) {
			apiKey, err := queries.GetAPIKey(c.Request().Context(), key)
			if err != nil {
				slog.Error("api key not found", slog.String("key", key), slog.Any("error", err))
				return false, err
			}
			c.Set(handlers.APIKeyCtxKey, &apiKey)
			return apiKey.Status == 0, nil
		},
	}))
	v1.POST("/chat/completions", openaiApiHandler.CreateChatCompletion)
	v1.POST("/audio/transcriptions", openaiApiHandler.CreateTranscription)
	v1.POST("/messages", anthropicApiHandler.CreateMessages)

	console := e.Group("/console", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "cookie:session",
		Validator: func(sessionID string, c echo.Context) (bool, error) {
			user, err := queries.GetUserBySessionID(c.Request().Context(), pgtype.Text{String: sessionID, Valid: true})
			if err != nil {
				slog.Error("get user by session id error", slog.Any("error", err))
				return !strings.HasPrefix(c.Request().URL.Path, "/console/"), err
			}
			c.Set("user", &user)
			return true, nil
		},
	}))
	console.GET("/home", consoleHandler.HomePage)
	console.POST("/home", consoleHandler.HomePage)
	console.GET("/create_key_page", consoleHandler.CreateKeyPage)
	console.POST("/generate_key", consoleHandler.GenerateKey)

	user := e.Group("/user")
	user.GET("/login", userHandler.LoginPage)
	user.POST("/login", userHandler.DoLogin)
}
