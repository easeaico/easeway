package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/easeaico/easeway/internal/config"
	"github.com/easeaico/easeway/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/liushuangls/go-anthropic/v2"
)

type AnthropicApiHandler struct {
	client  *anthropic.Client
	queries *store.Queries
}

func NewAthropicApiHandler(cfg *config.Config, queries *store.Queries) *AnthropicApiHandler {
	return &AnthropicApiHandler{
		client:  anthropic.NewClient(cfg.Anthropic.ApiKey),
		queries: queries,
	}
}

func (a *AnthropicApiHandler) CreateMessages(c echo.Context) error {
	startTime := time.Now().UnixMilli()

	ctx := c.Request().Context()
	key := c.Get(APIKeyCtxKey).(*store.ApiKey)

	var (
		promptTokens     = 0
		completionTokens = 0
		err              error
	)

	req := anthropic.MessagesRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	var resp anthropic.MessagesResponse
	if req.Stream {
		resp, err = a.client.CreateMessages(ctx, req)
	} else {
		resp, err = a.client.CreateMessagesStream(ctx, anthropic.MessagesStreamRequest{
			MessagesRequest: req,
			OnError: func(anthropic.ErrorResponse) {
			},
			OnPing: func(anthropic.MessagesEventPingData) {
			},
			OnMessageStart: func(anthropic.MessagesEventMessageStartData) {
			},
			OnContentBlockStart: func(anthropic.MessagesEventContentBlockStartData) {
			},
			OnContentBlockDelta: func(anthropic.MessagesEventContentBlockDeltaData) {
			},
			OnContentBlockStop: func(anthropic.MessagesEventContentBlockStopData, anthropic.MessageContent) {
			},
			OnMessageDelta: func(anthropic.MessagesEventMessageDeltaData) {
			},
			OnMessageStop: func(anthropic.MessagesEventMessageStopData) {
			},
		})
	}
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			slog.Error("Messages stream error, type: %s, message: %s", e.Type, e.Message)
		} else {
			slog.Error("do chat completion error", slog.Any("error", err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	promptTokens = resp.Usage.InputTokens
	completionTokens = resp.Usage.OutputTokens

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

	return c.JSON(http.StatusOK, resp)
}
