package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/easeaico/easeway/internal/config"
	"github.com/easeaico/easeway/internal/handlers"
	"github.com/easeaico/easeway/internal/spi"
	"github.com/easeaico/easeway/internal/spi/googleapi"
	"github.com/easeaico/easeway/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "f", "config.yaml", "config file path")
	flag.Parse()

	conf := config.NewConfig(confFile)

	ctx := context.Background()

	spis := spi.NewSPIRegistry(conf)
	gemini := googleapi.NewGenerativeAIClient(ctx, conf)
	spis.AddModelSPI("gemini-1.0-pro", gemini)

	pool := store.NewDBTX(ctx, conf)
	defer pool.Close()
	queries := store.New(pool)

	homeHandler := handlers.NewHomeHandler()
	apiHandler := handlers.NewAPISvcHandler(spis, queries)
	keyHandler := handlers.NewAPIKeyHandler(queries)
	userHandler := handlers.NewUserHandler(conf, queries)

	e := echo.New()
	e.Static("/assets", "assets")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/", homeHandler.HomePage)
	e.POST("/generate_key", keyHandler.GenerateNewKey)

	v1 := e.Group("/v1", middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		ctx := c.Request().Context()
		apiKey, err := queries.GetAPIKey(ctx, key)
		if err != nil {
			slog.Error("api key not found", slog.String("key", key), slog.Any("error", err))
			return false, err
		}

		c.Set(handlers.APIKeyCtxKey, &apiKey)
		return apiKey.Status == 0, nil
	}))
	v1.POST("/chat/completions", apiHandler.CreateChatCompletion)

	user := e.Group("/user")
	user.GET("/login", userHandler.LoginPage)
	user.POST("/do_login", userHandler.DoLogin)
	user.GET("/send_verification", userHandler.SendVerificationCode)

	err := e.Start(fmt.Sprintf("%s:%d", conf.Server.IP, conf.Server.Port))
	if err != nil {
		e.Logger.Fatal(err)
	}
}
