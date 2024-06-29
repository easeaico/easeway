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
	var confFile string
	flag.StringVar(&confFile, "f", "config.yaml", "config file path")
	flag.Parse()

	conf := config.NewConfig(confFile)
	ctx := context.Background()
	spis := spi.NewSPIRegistry(ctx, conf)

	pool := store.NewDBTX(ctx, conf)
	defer pool.Close()
	queries := store.New(pool)

	homeHandler := handlers.NewHomeHandler(queries)
	apiHandler := handlers.NewAPISvcHandler(spis, queries)
	userHandler := handlers.NewUserHandler(conf, queries)
	consoleHandler := handlers.NewConsoleHandler(queries)
	supportHandler := handlers.NewSupportHandler()
	memberHandler := handlers.NewMemberHandler()

	e := echo.New()
	e.File("/favicon.ico", "assets/images/favicon.ico")
	e.Static("/assets", "assets")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", homeHandler.HomePage)
	e.GET("/support", supportHandler.HomePage)
	e.GET("/member", memberHandler.HomePage)

	v1 := e.Group("/v1", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Validator: func(key string, c echo.Context) (bool, error) {
			ctx := c.Request().Context()
			apiKey, err := queries.GetAPIKey(ctx, key)
			if err != nil {
				slog.Error("api key not found", slog.String("key", key), slog.Any("error", err))
				return false, err
			}

			c.Set(handlers.APIKeyCtxKey, &apiKey)
			return apiKey.Status == 0, nil
		},
	}))
	v1.POST("/chat/completions", apiHandler.CreateChatCompletion)
	v1.POST("/audio/transcriptions", apiHandler.CreateTranslation)
	v1.POST("/audio/speech", apiHandler.CreateSpeech)

	console := e.Group("/console", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "cookie:session",
		Validator: func(sessionID string, c echo.Context) (bool, error) {
			user, err := queries.GetUserBySessionID(c.Request().Context(), pgtype.Text{
				String: sessionID,
				Valid:  true,
			})
			if err != nil {
				slog.Error("get user by session id error", slog.Any("error", err))
				if strings.HasPrefix(c.Request().URL.Path, "/console/") {
					return false, err
				} else {
					return true, err
				}
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

	err := e.Start(fmt.Sprintf("%s:%d", conf.Server.IP, conf.Server.Port))
	if err != nil {
		e.Logger.Fatal(err)
	}
}
