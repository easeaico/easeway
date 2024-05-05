package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/easeaico/easeway/internal/spi"
	"github.com/easeaico/easeway/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/pkoukk/tiktoken-go"
	tiktoken_loader "github.com/pkoukk/tiktoken-go-loader"
	"github.com/sashabaranov/go-openai"
)

const APIKeyCtxKey = "APIKey"

type APISvcHandler struct {
	spis    *spi.SPIRegistry
	queries *store.Queries
	tke     *tiktoken.Tiktoken
}

func NewAPISvcHandler(spis *spi.SPIRegistry, queries *store.Queries) *APISvcHandler {
	tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())
	tke, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		panic(err)
	}

	return &APISvcHandler{
		spis:    spis,
		queries: queries,
		tke:     tke,
	}
}

func (a *APISvcHandler) CreateChatCompletion(c echo.Context) error {
	req := openai.ChatCompletionRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	mspi := a.spis.LoadByModel(req.Model)
	ctx := c.Request().Context()
	r := c.Response()

	key := c.Get(APIKeyCtxKey).(*store.ApiKey)

	if req.Stream {
		return a.doChatCompletionStream(ctx, &req, mspi, key, r)
	} else {
		return a.doChatCompletion(ctx, &req, mspi, key, r)
	}
}

func (a *APISvcHandler) doChatCompletion(ctx context.Context, req *openai.ChatCompletionRequest, mspi spi.ModelSPI, key *store.ApiKey, r *echo.Response) error {
	startTime := time.Now().UnixMilli()

	resp, err := mspi.CreateChatCompletion(ctx, req)
	if err != nil {
		return err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	fmt.Fprintf(r, "data: %s\n\n", string(data))
	r.Flush()

	_, err = a.queries.CreateOutcome(ctx, store.CreateOutcomeParams{
		UserID:           1,
		KeyID:            key.ID,
		PromptTokens:     int32(resp.Usage.PromptTokens),
		CompletionTokens: int32(resp.Usage.CompletionTokens),
		TotalTokens:      int32(resp.Usage.TotalTokens),
		Rt:               int32(time.Now().UnixMilli() - startTime),
		FeeRate:          10,
		Cost:             int32(float32(resp.Usage.TotalTokens) * 1.1),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (a *APISvcHandler) doChatCompletionStream(ctx context.Context, req *openai.ChatCompletionRequest, mspi spi.ModelSPI, key *store.ApiKey, r *echo.Response) error {
	startTime := time.Now().UnixMilli()

	stream, err := mspi.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}

	r.Header().Set(echo.HeaderContentType, "text/event-stream")
	r.WriteHeader(http.StatusOK)

	promptTokens := 0
	for _, msg := range req.Messages {
		promptTokens += len(a.tke.Encode(msg.Content, nil, nil))
	}

	completionTokens := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			resp, err := stream.Recv()
			if err == io.EOF {
				totalTokens := promptTokens + completionTokens
				_, err = a.queries.CreateOutcome(ctx, store.CreateOutcomeParams{
					UserID:           1,
					KeyID:            key.ID,
					PromptTokens:     int32(promptTokens),
					CompletionTokens: int32(completionTokens),
					TotalTokens:      int32(totalTokens),
					Rt:               int32(time.Now().UnixMilli() - startTime),
					FeeRate:          10,
					Cost:             int32(float32(totalTokens) * 1.1),
				})
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				return nil
			}

			if err != nil {
				return err
			}

			completionTokens += len(a.tke.Encode(resp.Choices[0].Delta.Content, nil, nil))

			data, err := json.Marshal(resp)
			if err != nil {
				return err
			}

			fmt.Fprintf(r, "data: %s\n\n", string(data))
			r.Flush()
		}
	}
}
