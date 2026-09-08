
package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ollamaEmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type ollamaEmbeddingResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

const (
	ollamaURL   = "http://localhost:11434/api/embed"
	embeddingModel = "nomic-embed-text"
)

func VectorEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	embeddings := make([][]float32, 0, len(texts))

	for _, text := range texts {
		reqBody := ollamaEmbeddingRequest{
			Model: embeddingModel,
			Input: text,
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal embedding request: %w", err)
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			ollamaURL,
			bytes.NewReader(body),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create embedding request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to Ollama: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()

			return nil, fmt.Errorf(
				"Ollama returned status %d",
				resp.StatusCode,
			)
		}

		var result ollamaEmbeddingResponse

		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf(
				"failed to decode Ollama response: %w",
				err,
			)
		}

		if len(result.Embeddings) == 0 {
			return nil, fmt.Errorf("Ollama returned no embeddings")
		}

		embeddings = append(embeddings, result.Embeddings[0])
	}

	return embeddings, nil
}
