package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

func (a *APISvcHandler) CreateTranscription(c echo.Context) error {
	req := openai.AudioRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	slog.Info("create transcription", slog.Any("req", req))

	spi := a.spis.LoadByAsrModel(req.Model)
	if spi == nil {
		msg := fmt.Sprintf("unknown model name: %s", req.Model)
		slog.Error(msg, slog.String("model", req.Model))
		return c.String(http.StatusBadRequest, msg)
	}

	resp, err := spi.CreateTranscription(c.Request().Context(), &req)
	if err != nil {
		slog.Error("create transcription error", slog.String("model", req.Model))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *APISvcHandler) CreateSpeech(c echo.Context) error {
	req := openai.CreateSpeechRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	spi := a.spis.LoadByTtsModel(string(req.Model))
	if spi == nil {
		msg := fmt.Sprintf("unknown model name: %s", req.Model)
		slog.Error(msg, slog.String("model", string(req.Model)))
		return c.String(http.StatusBadRequest, msg)
	}

	ctx := c.Request().Context()
	resp, err := spi.CreateSpeech(ctx, &req)
	if err != nil {
		slog.Error("create transcription error", slog.String("model", string(req.Model)))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *APISvcHandler) CreateChatCompletion(c echo.Context) error {
	startTime := time.Now().UnixMilli()

	req := openai.ChatCompletionRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	mspi := a.spis.LoadByLlmModel(req.Model)
	if mspi == nil {
		msg := fmt.Sprintf("unknown model name: %s", req.Model)
		slog.Error(msg, slog.String("model", req.Model))
		return c.String(http.StatusBadRequest, msg)
	}

	ctx := c.Request().Context()
	r := c.Response()

	key := c.Get(APIKeyCtxKey).(*store.ApiKey)

	var (
		promptTokens     = 0
		completionTokens = 0
		err              error
	)

	if req.Stream {
		err = a.doChatCompletionStream(ctx, &req, mspi, &promptTokens, &completionTokens, r)
	} else {
		err = a.doChatCompletion(ctx, &req, mspi, &promptTokens, &completionTokens, c)
	}

	if err != nil {
		slog.Error("do chat completion error", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	totalTokens := promptTokens + completionTokens
	_, err = a.queries.CreateOutcome(ctx, store.CreateOutcomeParams{
		KeyID:            key.ID,
		UserID:           key.UserID,
		ModelName:        req.Model,
		PromptTokens:     int32(promptTokens),
		CompletionTokens: int32(completionTokens),
		TotalTokens:      int32(totalTokens),
		Rt:               int32(time.Now().UnixMilli() - startTime),
		FeeRate:          10,
		Cost:             int32(float32(totalTokens) * 1.1),
	})
	if err != nil {
		slog.Error("create outcome error", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (a *APISvcHandler) doChatCompletion(ctx context.Context, req *openai.ChatCompletionRequest, mspi spi.LlmSpi, promptTokens *int, completionTokens *int, c echo.Context) error {
	resp, err := mspi.CreateChatCompletion(ctx, req)
	if err != nil {
		return err
	}

	*promptTokens = resp.Usage.PromptTokens
	*completionTokens = resp.Usage.CompletionTokens
	return c.JSON(http.StatusOK, resp)
}

func (a *APISvcHandler) doChatCompletionStream(ctx context.Context, req *openai.ChatCompletionRequest, mspi spi.LlmSpi, promptTokens *int, completionTokens *int, r *echo.Response) error {
	stream, err := mspi.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}

	r.Header().Set(echo.HeaderContentType, "text/event-stream")
	r.WriteHeader(http.StatusOK)

	for _, msg := range req.Messages {
		*promptTokens += len(a.tke.Encode(msg.Content, nil, nil))
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			resp, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}

			*completionTokens += len(a.tke.Encode(resp.Choices[0].Delta.Content, nil, nil))

			data, err := json.Marshal(resp)
			if err != nil {
				return err
			}

			fmt.Fprintf(r, "data: %s\n\n", string(data))
			r.Flush()
		}
	}
}
