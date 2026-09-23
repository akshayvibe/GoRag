package generation

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type Generator struct {
	client *openai.Client
	model  string
}

func NewGenerator(apiKey string) *Generator {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.x.ai/v1"),
	)

	return &Generator{
		client: &client,
		model:  "grok-4.7",
	}
}

func (g *Generator) Generate(
	ctx context.Context,
	query string,
	contextChunks []string,
) (string, error) {

	contextText := strings.Join(contextChunks, "\n\n---\n\n")

	prompt := fmt.Sprintf(`
You are a helpful RAG assistant.

Answer the user's question using ONLY the provided context.

If the answer cannot be found in the context, say:
"I don't know based on the provided documents."

Do not make up information.

CONTEXT:
%s

QUESTION:
%s

ANSWER:
`, contextText, query)

	response, err := g.client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model: openai.ChatModel(g.model),
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(prompt),
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("grok generation failed: %w", err)
	}

	return response.OutputText(), nil
}