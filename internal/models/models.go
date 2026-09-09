package models

type Chunk struct {
	ID       string  `json:"id"`
	Content  string  `json:"content"`
	Tokens   int     `json:"tokens"`
	Vector   []float32 `json:"vector"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	DocumentID string `json:"document_id"`
	FileName   string `json:"file_name"`
	// PageNumber int    `json:"page_number"`
	// ChunkIndex int    `json:"chunk_index"`
}

type ChunkPayload struct {
	Content  string   `json:"content"`
	Tokens   int      `json:"tokens"`
	Metadata Metadata `json:"metadata"`
}

type VectorPoint struct {
	ID      string
	Vector  []float32
	Payload ChunkPayload
}
type PdfContent struct {
	Text     string
	Metadata Metadata
}