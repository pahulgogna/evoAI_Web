package types

import "context"

type LLMInterface interface {
	StreamMessage(ctx context.Context, messages *[]Message, model string) (<-chan string, error)
}
