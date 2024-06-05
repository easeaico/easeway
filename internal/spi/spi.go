package spi

import (
	"context"

	"github.com/easeaico/easeway/internal/config"
	"github.com/easeaico/easeway/internal/spi/google"
	"github.com/easeaico/easeway/internal/spi/groq"
	openaispi "github.com/easeaico/easeway/internal/spi/openai"
	"github.com/sashabaranov/go-openai"
)

type ModelSPI interface {
	CreateChatCompletionStream(
		ctx context.Context,
		request *openai.ChatCompletionRequest,
	) (stream *openai.ChatCompletionStream, err error)

	CreateChatCompletion(
		ctx context.Context,
		request *openai.ChatCompletionRequest,
	) (response *openai.ChatCompletionResponse, err error)
}

type SPIRegistry struct {
	conf      *config.Config
	providers map[string]ModelSPI
}

func NewSPIRegistry(ctx context.Context, conf *config.Config) *SPIRegistry {
	groqSPI := groq.NewGroqClient(conf)
	geminiSPI := google.NewGenerativeAIClient(ctx, conf)
	openaiSPI := openaispi.NewOpenAIClient(conf)

	providers := map[string]ModelSPI{
		"gpt-4-32k-0613":         openaiSPI,
		"gpt-4-32k-0314":         openaiSPI,
		"gpt-4-32k":              openaiSPI,
		"gpt-4-0613":             openaiSPI,
		"gpt-4-0314":             openaiSPI,
		"gpt-4o":                 openaiSPI,
		"gpt-4o-2024-05-13":      openaiSPI,
		"gpt-4-turbo":            openaiSPI,
		"gpt-4-turbo-2024-04-09": openaiSPI,
		"gpt-4-0125-preview":     openaiSPI,
		"gpt-4-1106-preview":     openaiSPI,
		"gpt-4-turbo-preview":    openaiSPI,
		"gpt-4-vision-preview":   openaiSPI,
		"gpt-4":                  openaiSPI,
		"gpt-3.5-turbo-0125":     openaiSPI,
		"gpt-3.5-turbo-1106":     openaiSPI,
		"gpt-3.5-turbo-0613":     openaiSPI,
		"gpt-3.5-turbo-0301":     openaiSPI,
		"gpt-3.5-turbo-16k":      openaiSPI,
		"gpt-3.5-turbo-16k-0613": openaiSPI,
		"gpt-3.5-turbo":          openaiSPI,
		"gpt-3.5-turbo-instruct": openaiSPI,
		"gemini-1.0-pro":         geminiSPI,
		"gemma-7b-it":            groqSPI,
		"llama3-8b-8192":         groqSPI,
		"llama3-70b-8192":        groqSPI,
		"mixtral-8x7b-32768":     groqSPI,
	}
	return &SPIRegistry{
		conf:      conf,
		providers: providers,
	}
}

func (r SPIRegistry) LoadByModel(model string) ModelSPI {
	return r.providers[model]
}

func (r *SPIRegistry) AddModelSPI(model string, spi ModelSPI) {
	r.providers[model] = spi
}
