package models

type Chunk struct {
    ID       string
    Content  string
    Metadata Metadata
}

type Metadata struct {
    DocumentID string
    FileName   string
    PageNumber int
    ChunkIndex int
}