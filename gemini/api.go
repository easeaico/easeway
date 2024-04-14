package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	gl "cloud.google.com/go/ai/generativelanguage/apiv1beta"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1beta/generativelanguagepb"
	"github.com/easeaico/easeway/config"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

type GenerativeAIClient struct {
	client *gl.GenerativeClient
}

func NewGenerativeAIClient(conf *config.Config) *GenerativeAIClient {
	ctx := context.Background()
	client, err := gl.NewGenerativeRESTClient(ctx, option.WithAPIKey(conf.Gemini.ApiKey))
	if err != nil {
		log.Fatal(err)
	}
	return &GenerativeAIClient{
		client: client,
	}
}

func (g *GenerativeAIClient) Handle(c echo.Context) error {
	modelName := c.Param("modelName")
	ctx := c.Request().Context()

	req := &pb.GenerateContentRequest{}
	if err := c.Bind(req); err != nil {
		return fmt.Errorf("parase content request error: %w", err)
	}
	req.Model = modelName
	stream, err := g.client.StreamGenerateContent(c.Request().Context(), req)
	if err != nil {
		return fmt.Errorf("stream generate content error: %w", err)
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

func (g *GenerativeAIClient) Close() error {
	return g.client.Close()
}
