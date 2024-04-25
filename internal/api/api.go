package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/easeaico/easeway/internal/config"
	"github.com/easeaico/easeway/internal/spi"
	"github.com/easeaico/easeway/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

const APIKeyCtxKey = "APIKey"

type ApiSvc struct {
	conf *config.Config
	spis *spi.SPIRegistry
}

func NewApiSvc(conf *config.Config, spis *spi.SPIRegistry) *ApiSvc {
	return &ApiSvc{
		conf: conf,
		spis: spis,
	}
}

func (a *ApiSvc) CreateChatCompletion(c echo.Context) error {
	req := openai.ChatCompletionRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	mspi := a.spis.LoadByModel(req.Model)
	ctx := c.Request().Context()
	_ = c.Get(APIKeyCtxKey).(*store.ApiKey)
	r := c.Response()
	if req.Stream {
		stream, err := mspi.CreateChatCompletionStream(ctx, &req)
		if err != nil {
			return err
		}

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
	} else {
		resp, err := mspi.CreateChatCompletion(ctx, &req)
		if err != nil {
			return err
		}

		data, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		fmt.Fprintf(r, "data: %s\n\n", string(data))
		r.Flush()
		return nil
	}
}
