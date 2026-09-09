package db

import "github.com/qdrant/go-client/qdrant"

type VectorStore interface {
	CreateCollection() error
	CollectionExists() (bool, error)

	Upsert(points []*qdrant.PointStruct) error

	Search(
		vector []float32,
		limit uint64,
	) ([]*qdrant.ScoredPoint, error)

	Delete(ids []uint64) error
}