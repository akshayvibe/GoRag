package retrieval

import (
	"context"
	"fmt"

	"github.com/akshayvibe/GoRag/internal/db/vectorstore"
	"github.com/akshayvibe/GoRag/internal/indexing/embedding"
	"github.com/akshayvibe/GoRag/internal/models"
)

type Retriever struct {
	Store *vectorstore.QdrantStore
}

func NewRetriever(store *vectorstore.QdrantStore) *Retriever {
	return &Retriever{
		Store: store,
	}
}

func (r *Retriever) Retrieve(
	ctx context.Context,
	query string,
	topK uint64,
) ([]models.VectorPoint, error) {

	// 1. Convert the user's query into a vector
	vectors, err := embedding.VectorEmbedding(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	if len(vectors) == 0 {
		return nil, fmt.Errorf("no embedding returned for query")
	}

	// 2. Search Qdrant
	results, err := r.Store.Search(ctx, vectors[0], topK)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %w", err)
	}

	// 3. Convert Qdrant results into our application model
	points := make([]models.VectorPoint, 0, len(results))

	for _, result := range results {
		point := models.VectorPoint{
			ID: fmt.Sprintf("%d", result.Id.GetNum()),
		}

		if result.Payload != nil {
			if textValue, ok := result.Payload["text"]; ok {
				point.Payload.Content = textValue.GetStringValue()
			}

			if sourceValue, ok := result.Payload["source"]; ok {
				point.Payload.Metadata.FileName = sourceValue.GetStringValue()
			}
		}

		points = append(points, point)
	}

	return points, nil
}