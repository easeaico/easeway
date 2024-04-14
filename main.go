package main

import (
	"flag"
	"fmt"

	"github.com/easeaico/easeway/config"
	"github.com/easeaico/easeway/gemini"
	"github.com/easeaico/easeway/openai"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "f", "config.yaml", "配置文件路径")
	flag.Parse()

	conf := config.NewConfig(confFile)
	openaiClient := openai.NewOpenAIClient(&conf)
	geminiClient := gemini.NewGenerativeAIClient(&conf)
	defer geminiClient.Close()

	e := echo.New()
	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/api/openai/v1/chat/completions", openaiClient.Handle)
	e.POST("/api/openai/gemini/v1beta/models/:modelName\\:streamGenerateContent?key=:apiKey", geminiClient.Handle)

	err := e.Start(fmt.Sprintf("%s:%d", conf.Server.IP, conf.Server.Port))
	if err != nil {
		e.Logger.Fatal(err)
	}
}
