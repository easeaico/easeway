package groq

import (
	"context"

	"github.com/easeaico/easeway/internal/config"
	"github.com/sashabaranov/go-openai"
)

type GroqClient struct {
	client *openai.Client
}

func NewGroqClient(conf *config.Config) *GroqClient {
	cfg := openai.DefaultConfig(conf.Groq.ApiKey)
	cfg.BaseURL = "https://api.groq.com/openai/v1"
	client := openai.NewClientWithConfig(cfg)
	return &GroqClient{
		client: client,
	}
}

func (o *GroqClient) CreateChatCompletionStream(
	ctx context.Context,
	request *openai.ChatCompletionRequest,
) (stream *openai.ChatCompletionStream, err error) {
	return o.client.CreateChatCompletionStream(ctx, *request)
}

func (o *GroqClient) CreateChatCompletion(
	ctx context.Context,
	request *openai.ChatCompletionRequest,
) (response *openai.ChatCompletionResponse, err error) {
	resp, err := o.client.CreateChatCompletion(ctx, *request)
	return &resp, err
}
