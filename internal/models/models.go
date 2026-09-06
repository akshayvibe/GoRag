package models

type Chunk struct {
    ID       string  `json:"id"`
    Content  string  `json:"content"`
    Tokens   int     `json:"tokens"`
    Metadata Metadata `json:"metadata"`
}

type Metadata struct {
    DocumentID string `json:"document_id"`
    FileName   string `json:"file_name"`
    // PageNumber int    `json:"page_number"`
    // ChunkIndex int    `json:"chunk_index"`
}