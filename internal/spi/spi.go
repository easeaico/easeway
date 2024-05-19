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
		openai.GPT432K0613:           openaiAPI,
		openai.GPT432K0314:           openaiAPI,
		openai.GPT432K:               openaiAPI,
		openai.GPT40613:              openaiAPI,
		openai.GPT40314:              openaiAPI,
		openai.GPT4o:                 openaiAPI,
		openai.GPT4o20240513:         openaiAPI,
		openai.GPT4Turbo:             openaiAPI,
		openai.GPT4Turbo20240409:     openaiAPI,
		openai.GPT4Turbo0125:         openaiAPI,
		openai.GPT4Turbo1106:         openaiAPI,
		openai.GPT4TurboPreview:      openaiAPI,
		openai.GPT4VisionPreview:     openaiAPI,
		openai.GPT4:                  openaiAPI,
		openai.GPT3Dot5Turbo0125:     openaiAPI,
		openai.GPT3Dot5Turbo1106:     openaiAPI,
		openai.GPT3Dot5Turbo0613:     openaiAPI,
		openai.GPT3Dot5Turbo0301:     openaiAPI,
		openai.GPT3Dot5Turbo16K:      openaiAPI,
		openai.GPT3Dot5Turbo16K0613:  openaiAPI,
		openai.GPT3Dot5Turbo:         openaiAPI,
		openai.GPT3Dot5TurboInstruct: openaiAPI,
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
