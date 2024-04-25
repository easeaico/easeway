package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"

	"github.com/easeaico/easeway/internal/api"
	"github.com/easeaico/easeway/internal/config"
	"github.com/easeaico/easeway/internal/spi"
	"github.com/easeaico/easeway/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "f", "config.yaml", "配置文件路径")
	flag.Parse()

	conf := config.NewConfig(confFile)
	spis := spi.NewSPIRegistry(conf)
	apiSvc := api.NewApiSvc(conf, spis)

	db, err := sql.Open("sqlite3", conf.DBFile)
	if err != nil {
		slog.Error("sqlite db open error", slog.Any("error", err))
		return
	}
	queries := store.New(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		ctx := c.Request().Context()
		apiKey, err := queries.GetAPIKey(ctx, key)
		if err != nil {
			slog.Error("api key not found", slog.String("key", key), slog.Any("error", err))
			return false, err
		}

		c.Set(api.APIKeyCtxKey, key)
		return apiKey.Status == 0, nil
	}))
	e.POST("/v1/chat/completions", apiSvc.CreateChatCompletion)

	err = e.Start(fmt.Sprintf("%s:%d", conf.Server.IP, conf.Server.Port))
	if err != nil {
		e.Logger.Fatal(err)
	}
}
