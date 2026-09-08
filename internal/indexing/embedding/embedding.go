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

func VectorEmbedding(texts []string) ([][]float32, error) {
	ctx := context.Background()

	embeddings := make([][]float32, 0, len(texts))

	for _, text := range texts {
		reqBody := ollamaEmbeddingRequest{
			Model: "nomic-embed-text",
			Input: text,
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			"http://localhost:11434/api/embed",
			bytes.NewBuffer(body),
		)
		if err != nil {
			return nil, err
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
			return nil, err
		}

		if len(result.Embeddings) == 0 {
			return nil, fmt.Errorf("Ollama returned no embeddings")
		}

		embeddings = append(embeddings, result.Embeddings[0])
	}

	return embeddings, nil
}

