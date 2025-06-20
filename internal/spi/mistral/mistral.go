package mistral

import (
	"context"

	"github.com/easeaico/easeway/internal/config"
	mistralgo "github.com/gage-technologies/mistral-go"
	"github.com/sashabaranov/go-openai"
)

type MistralClient struct {
	client *mistralgo.MistralClient
}

func NewMistralClient(conf *config.Config) *MistralClient {
	client := mistralgo.NewMistralClientDefault(conf.Mistral.ApiKey)
	return &MistralClient{
		client: client,
	}
}

func (o *MistralClient) CreateChatCompletionStream(
	ctx context.Context,
	request *openai.ChatCompletionRequest,
) (stream *openai.ChatCompletionStream, err error) {
	return nil, nil
}

func (o *MistralClient) CreateChatCompletion(
	ctx context.Context,
	request *openai.ChatCompletionRequest,
) (response *openai.ChatCompletionResponse, err error) {
	return nil, nil
}
