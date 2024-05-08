package spi

import (
	"context"

	"github.com/easeaico/easeway/internal/config"
	"github.com/easeaico/easeway/internal/spi/openaiapi"
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

func NewSPIRegistry(conf *config.Config) *SPIRegistry {
	openaiAPI := openaiapi.NewOpenAIClient(conf)

	providers := map[string]ModelSPI{
		openai.GPT4TurboPreview: openaiAPI,
		openai.GPT4Turbo0125:    openaiAPI,
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
