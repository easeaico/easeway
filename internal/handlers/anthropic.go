package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/liushuangls/go-anthropic/v2"
)

type AnthropicApiHandler struct {
	client *anthropic.Client
}

func NewAthropicApiHandler() *AnthropicApiHandler {
	return &AnthropicApiHandler{
		client: anthropic.NewClient(""),
	}
}

func (a *AnthropicApiHandler) CreateMessages(c echo.Context) error {
	req := anthropic.MessagesRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	var (
		resp anthropic.MessagesResponse
		err  error
	)
	if req.Stream {
		resp, err = a.client.CreateMessages(c.Request().Context(), req)
	} else {
		resp, err = a.client.CreateMessagesStream(c.Request().Context(), anthropic.MessagesStreamRequest{
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
			fmt.Printf("Messages stream error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages stream error: %v\n", err)
		}
	}

	return c.JSON(http.StatusOK, resp)
}
