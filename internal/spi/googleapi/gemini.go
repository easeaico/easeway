package googleapi

import (
	"context"
	"encoding/base64"
	"log"
	"strings"

	"github.com/easeaico/easeway/internal/config"
	"github.com/google/generative-ai-go/genai"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
)

type GenerativeAIClient struct {
	client *genai.Client
}

func NewGenerativeAIClient(ctx context.Context, conf *config.Config) *GenerativeAIClient {
	client, err := genai.NewClient(ctx, option.WithAPIKey(conf.Gemini.ApiKey))
	if err != nil {
		log.Fatal(err)
	}

	return &GenerativeAIClient{
		client: client,
	}
}

func (g *GenerativeAIClient) CreateChatCompletionStream(
	ctx context.Context,
	req *openai.ChatCompletionRequest,
) (*openai.ChatCompletionStream, error) {
	model := g.client.GenerativeModel(req.Model)
	cs := model.StartChat()
	histories, err := convertHistories(req.Messages)
	if err != nil {
		return nil, err
	}
	cs.History = histories

	_ = cs.SendMessageStream(ctx)
	stream := &openai.ChatCompletionStream{}
	return stream, nil
}

func (g *GenerativeAIClient) CreateChatCompletion(
	ctx context.Context,
	req *openai.ChatCompletionRequest,
) (*openai.ChatCompletionResponse, error) {
	model := g.client.GenerativeModel(req.Model)
	cs := model.StartChat()

	histories, err := convertHistories(req.Messages)
	if err != nil {
		return nil, err
	}
	cs.History = histories

	resp, err := cs.SendMessage(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var (
		completionTokens int
		choices          []openai.ChatCompletionChoice
	)
	for _, c := range resp.Candidates {
		role := c.Content.Role
		if role == "model" {
			role = openai.ChatMessageRoleAssistant
		}

		content := c.Content.Parts[0].(genai.Text)
		msg := openai.ChatCompletionMessage{
			Role:    role,
			Content: string(content),
		}
		choice := openai.ChatCompletionChoice{
			Index:        0,
			Message:      msg,
			FinishReason: "",
			LogProbs:     &openai.LogProbs{},
		}
		choices = append(choices, choice)
	}

	sresp := openai.ChatCompletionResponse{
		Model:   req.Model,
		Choices: choices,
		Usage: openai.Usage{
			CompletionTokens: completionTokens,
		},
	}
	return &sresp, nil
}

func (g *GenerativeAIClient) Close() error {
	return g.client.Close()
}

func convertHistories(messages []openai.ChatCompletionMessage) ([]*genai.Content, error) {
	var histories []*genai.Content
	for _, msg := range messages {
		role := msg.Role
		if role == openai.ChatMessageRoleAssistant {
			role = "model"
		}

		history := genai.Content{
			Role: role,
		}
		if len(msg.Content) > 0 {
			history.Parts = []genai.Part{
				genai.Text(msg.Content),
			}
		} else if len(msg.MultiContent) > 0 {
			var parts []genai.Part
			for _, mc := range msg.MultiContent {
				var part genai.Part
				switch mc.Type {
				case openai.ChatMessagePartTypeText:
					part = genai.Text(mc.Text)
				case openai.ChatMessagePartTypeImageURL:
					data := mc.ImageURL.URL
					datas := strings.Split(data, ";")
					imgData := datas[1][7:]
					img, err := base64.StdEncoding.DecodeString(imgData)
					if err != nil {
						return nil, err
					}
					part = genai.ImageData(datas[0][12:], img)
				}

				parts = append(parts, part)
			}
			history.Parts = parts
		}

		histories = append(histories, &history)
	}

	return histories, nil
}
