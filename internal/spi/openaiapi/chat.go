package openaiapi

import (
	"context"

	"github.com/easeaico/easeway/internal/config"
	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(conf *config.Config) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(conf.OpenAI.ApiKey),
	}
}

func (o *OpenAIClient) CreateChatCompletionStream(
	ctx context.Context,
	request *openai.ChatCompletionRequest,
) (stream *openai.ChatCompletionStream, err error) {
	return o.client.CreateChatCompletionStream(ctx, *request)
}

func (o *OpenAIClient) CreateChatCompletion(
	ctx context.Context,
	request *openai.ChatCompletionRequest,
) (response *openai.ChatCompletionResponse, err error) {
	resp, err := o.client.CreateChatCompletion(ctx, *request)
	return &resp, err
}
