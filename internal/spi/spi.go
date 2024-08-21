package spi

import (
	"context"

	"github.com/easeaico/easeway/internal/config"
	"github.com/easeaico/easeway/internal/spi/google"
	"github.com/easeaico/easeway/internal/spi/groq"
	"github.com/easeaico/easeway/internal/spi/mistral"
	openaispi "github.com/easeaico/easeway/internal/spi/openai"
	"github.com/sashabaranov/go-openai"
)

type LlmSpi interface {
	CreateChatCompletionStream(
		ctx context.Context,
		request *openai.ChatCompletionRequest,
	) (stream *openai.ChatCompletionStream, err error)

	CreateChatCompletion(
		ctx context.Context,
		request *openai.ChatCompletionRequest,
	) (response *openai.ChatCompletionResponse, err error)
}

type AsrSpi interface {
	CreateTranscription(
		ctx context.Context,
		request *openai.AudioRequest,
	) (response *openai.AudioResponse, err error)
}

type TtsSpi interface {
	CreateSpeech(
		ctx context.Context,
		request *openai.CreateSpeechRequest,
	) (response *openai.RawResponse, err error)
}

type SPIRegistry struct {
	conf         *config.Config
	llmProviders map[string]LlmSpi
	asrProviders map[string]AsrSpi
	ttsProviders map[string]TtsSpi
}

func NewSPIRegistry(ctx context.Context, conf *config.Config) *SPIRegistry {
	groqSPI := groq.NewGroqClient(conf)
	geminiSPI := google.NewGenerativeAIClient(ctx, conf)
	openaiSPI := openaispi.NewOpenAIClient(conf)
	mistralSPI := mistral.NewMistralClient(conf)

	asrProviders := map[string]AsrSpi{
		openai.Whisper1: openaiSPI,
	}

	ttsProviders := map[string]TtsSpi{
		string(openai.TTSModel1):      openaiSPI,
		string(openai.TTSModel1HD):    openaiSPI,
		string(openai.TTSModelCanary): openaiSPI,
	}

	llmProviders := map[string]LlmSpi{
		"gpt-4-32k-0613":          openaiSPI,
		"gpt-4-32k-0314":          openaiSPI,
		"gpt-4-32k":               openaiSPI,
		"gpt-4-0613":              openaiSPI,
		"gpt-4-0314":              openaiSPI,
		"gpt-4o":                  openaiSPI,
		"gpt-4o-2024-05-13":       openaiSPI,
		"gpt-4-turbo":             openaiSPI,
		"gpt-4-turbo-2024-04-09":  openaiSPI,
		"gpt-4-0125-preview":      openaiSPI,
		"gpt-4-1106-preview":      openaiSPI,
		"gpt-4-turbo-preview":     openaiSPI,
		"gpt-4-vision-preview":    openaiSPI,
		"gpt-4":                   openaiSPI,
		"gpt-3.5-turbo-0125":      openaiSPI,
		"gpt-3.5-turbo-1106":      openaiSPI,
		"gpt-3.5-turbo-0613":      openaiSPI,
		"gpt-3.5-turbo-0301":      openaiSPI,
		"gpt-3.5-turbo-16k":       openaiSPI,
		"gpt-3.5-turbo-16k-0613":  openaiSPI,
		"gpt-3.5-turbo":           openaiSPI,
		"gpt-3.5-turbo-instruct":  openaiSPI,
		"gemini-1.0-pro":          geminiSPI,
		"gemma-7b-it":             groqSPI,
		"llama3-8b-8192":          groqSPI,
		"llama3-70b-8192":         groqSPI,
		"llama-3.1-70b-versatile": groqSPI,
		"llama-3.1-8b-instant":    groqSPI,
		"mixtral-8x7b-32768":      groqSPI,
		"open-mixtral-7b":         mistralSPI,
		"open-mixtral-8x7b":       mistralSPI,
		"open-mixtral-8x22b":      mistralSPI,
		"mistral-small-latest":    mistralSPI,
		"mistral-medium-latest":   mistralSPI,
		"mistral-large-latest":    mistralSPI,
		"codestral-latest":        mistralSPI,
	}
	return &SPIRegistry{
		conf:         conf,
		llmProviders: llmProviders,
		asrProviders: asrProviders,
		ttsProviders: ttsProviders,
	}
}

func (r SPIRegistry) LoadByLlmModel(model string) LlmSpi {
	return r.llmProviders[model]
}

func (r SPIRegistry) LoadByAsrModel(model string) AsrSpi {
	return r.asrProviders[model]
}

func (r SPIRegistry) LoadByTtsModel(model string) TtsSpi {
	return r.ttsProviders[model]
}
