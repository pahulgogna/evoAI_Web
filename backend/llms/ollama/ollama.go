package ollama

import (
	"context"
	"log"

	"github.com/ollama/ollama/api"
	"github.com/pahulgogna/evoAI_Web/backend/types"
)

type Ollama struct {
}

func NewOllamaInterface() *Ollama {
	return &Ollama{}
}

func (o *Ollama) StreamMessage(ctx context.Context, messages *[]types.Message, model string) (<-chan string, error) {

	out := make(chan string, 1024)

	client, err := api.ClientFromEnvironment()
	if err != nil {
		close(out)
		return nil, err
	}
	

	go func(ctx context.Context, messages *[]types.Message, model string) {
		defer close(out)

		mappedMessages := []api.Message{}

		for _, m := range *messages {
			mappedMessages = append(mappedMessages, api.Message{
				Role:    string(m.Role),
				Content: m.Content,
			})
		}

		stream := true
		req := &api.ChatRequest{
			Model:    model,
			Stream:   &stream,
			Messages: mappedMessages,
		}

		respFunc := func(resp api.ChatResponse) error {
			out <- resp.Message.Content
			return nil
		}

		err := client.Chat(ctx, req, respFunc)
		if err != nil {
			log.Printf("error starting chat: %s", err.Error())
			return
		}
	}(ctx, messages, model)

	return out, nil
}


