
package vectorstore

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantStore struct {
	client         *qdrant.Client
	collectionName string
	vectorSize     uint64
}

// NewQdrantStore creates a new Qdrant vector store.
func NewQdrantStore(
	host string,
	port int,
	collectionName string,
	vectorSize uint64,
) (*QdrantStore, error) {

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to qdrant: %w", err)
	}

	return &QdrantStore{
		client:         client,
		collectionName: collectionName,
		vectorSize:     vectorSize,
	}, nil
}

// CreateCollection creates the Qdrant collection if it doesn't already exist.
func (s *QdrantStore) CreateCollection(ctx context.Context) error {

	exists, err := s.CollectionExists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check collection: %w", err)
	}

	if exists {
		return nil
	}

	err = s.client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: s.collectionName,

		VectorsConfig: qdrant.NewVectorsConfig(
			&qdrant.VectorParams{
				Size:     s.vectorSize,
				Distance: qdrant.Distance_Cosine,
			},
		),
	})

	if err != nil {
		return fmt.Errorf(
			"failed to create collection %q: %w",
			s.collectionName,
			err,
		)
	}

	return nil
}

// CollectionExists checks whether the collection exists.
func (s *QdrantStore) CollectionExists(
	ctx context.Context,
) (bool, error) {

	collections, err := s.client.ListCollections(ctx)
	if err != nil {
		return false, err
	}

	for _, collection := range collections {
		if collection == s.collectionName {
			return true, nil
		}
	}

	return false, nil
}

// Add inserts one vector into Qdrant.
func (s *QdrantStore) Add(
	ctx context.Context,
	id uint64,
	vector []float32,
	text string,
) error {

	if uint64(len(vector)) != s.vectorSize {
		return fmt.Errorf(
			"invalid vector size: expected %d, got %d",
			s.vectorSize,
			len(vector),
		)
	}

	_, err := s.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: s.collectionName,

		Points: []*qdrant.PointStruct{
			{
				Id: qdrant.NewIDNum(id),

				Vectors: qdrant.NewVectors(
					&qdrant.Vector{
						Data: vector,
					},
				),

				Payload: qdrant.NewValueMap(
					map[string]any{
						"text": text,
					},
				),
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to add vector: %w", err)
	}

	return nil
}

// AddDocument inserts a vector together with useful document metadata.
func (s *QdrantStore) AddDocument(
	ctx context.Context,
	id uint64,
	vector []float32,
	text string,
	source string,
) error {

	if uint64(len(vector)) != s.vectorSize {
		return fmt.Errorf(
			"invalid vector size: expected %d, got %d",
			s.vectorSize,
			len(vector),
		)
	}

	_, err := s.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: s.collectionName,

		Points: []*qdrant.PointStruct{
			{
				Id: qdrant.NewIDNum(id),

				Vectors: qdrant.NewVectors(
					&qdrant.Vector{
						Data: vector,
					},
				),

				Payload: qdrant.NewValueMap(
					map[string]any{
						"text":   text,
						"source": source,
					},
				),
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to add document: %w", err)
	}

	return nil
}

// AddDocuments inserts multiple vectors in one request.
func (s *QdrantStore) AddDocuments(
	ctx context.Context,
	ids []uint64,
	vectors [][]float32,
	texts []string,
	sources []string,
) error {

	if len(ids) != len(vectors) ||
		len(ids) != len(texts) ||
		len(ids) != len(sources) {

		return fmt.Errorf(
			"ids, vectors, texts and sources must have the same length",
		)
	}

	points := make([]*qdrant.PointStruct, 0, len(ids))

	for i := range ids {

		if uint64(len(vectors[i])) != s.vectorSize {
			return fmt.Errorf(
				"invalid vector size for point %d: expected %d, got %d",
				ids[i],
				s.vectorSize,
				len(vectors[i]),
			)
		}

		points = append(points, &qdrant.PointStruct{
			Id: qdrant.NewIDNum(ids[i]),

			Vectors: qdrant.NewVectors(
				&qdrant.Vector{
					Data: vectors[i],
				},
			),

			Payload: qdrant.NewValueMap(
				map[string]any{
					"text":   texts[i],
					"source": sources[i],
				},
			),
		})
	}

	_, err := s.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: s.collectionName,
		Points:         points,
	})

	if err != nil {
		return fmt.Errorf("failed to add documents: %w", err)
	}

	return nil
}

// Search finds the most similar vectors.
func (s *QdrantStore) Search(
	ctx context.Context,
	vector []float32,
	limit uint64,
) ([]*qdrant.ScoredPoint, error) {

	if uint64(len(vector)) != s.vectorSize {
		return nil, fmt.Errorf(
			"invalid query vector size: expected %d, got %d",
			s.vectorSize,
			len(vector),
		)
	}

	results, err := s.client.Query(
		ctx,
		&qdrant.QueryPoints{
			CollectionName: s.collectionName,

			Query: qdrant.NewQuery(
				&qdrant.VectorInput{
					{
						Variant: &qdrant.VectorInput_Vector{
							Vector: &qdrant.Vector{
								Data: vector,
							},
						},
					},
				},
			),

			Limit: &limit,

			WithPayload: qdrant.NewWithPayload(true),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	return results, nil
}

// Get retrieves points by their IDs.
func (s *QdrantStore) Get(
	ctx context.Context,
	ids []uint64,
) ([]*qdrant.RetrievedPoint, error) {

	pointIDs := make([]*qdrant.PointId, 0, len(ids))

	for _, id := range ids {
		pointIDs = append(pointIDs, qdrant.NewIDNum(id))
	}

	result, err := s.client.Get(
		ctx,
		&qdrant.GetPoints{
			CollectionName: s.collectionName,
			Ids:            pointIDs,

			WithPayload: qdrant.NewWithPayload(true),
			WithVectors: qdrant.NewWithVectors(false),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get points: %w", err)
	}

	return result, nil
}

// Delete deletes points by their IDs.
func (s *QdrantStore) Delete(
	ctx context.Context,
	ids []uint64,
) error {

	pointIDs := make([]*qdrant.PointId, 0, len(ids))

	for _, id := range ids {
		pointIDs = append(pointIDs, qdrant.NewIDNum(id))
	}

	_, err := s.client.Delete(
		ctx,
		&qdrant.DeletePoints{
			CollectionName: s.collectionName,

			Points: &qdrant.PointsSelector{
				PointsSelectorOneOf: &qdrant.PointsSelector_Points{
					Points: &qdrant.PointsIds{
						Ids: pointIDs,
					},
				},
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to delete points: %w", err)
	}

	return nil
}

// DeleteCollection permanently removes the collection.
func (s *QdrantStore) DeleteCollection(
	ctx context.Context,
) error {

	_, err := s.client.DeleteCollection(
		ctx,
		&qdrant.DeleteCollection{
			CollectionName: s.collectionName,
		},
	)

	if err != nil {
		return fmt.Errorf(
			"failed to delete collection %q: %w",
			s.collectionName,
			err,
		)
	}

	return nil
}

// Close closes the Qdrant client connection.
func (s *QdrantStore) Close() error {
	return s.client.Close()
}