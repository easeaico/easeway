package openai

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/easeaico/easeway/config"
	"github.com/labstack/echo/v4"
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

func (o *OpenAIClient) Handle(c echo.Context) error {
	req := openai.ChatCompletionRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	ctx := c.Request().Context()
	stream, err := o.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}

	r := c.Response()
	r.Header().Set(echo.HeaderContentType, "text/event-stream")
	r.WriteHeader(http.StatusOK)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			resp, err := stream.Recv()
			if err != nil {
				return err
			}

			data, err := json.Marshal(resp)
			if err != nil {
				return err
			}

			fmt.Fprintf(r, "data: %s\n\n", string(data))
			r.Flush()
		}
	}
}
