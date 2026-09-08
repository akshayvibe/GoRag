package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/akshayvibe/GoRag/internal/indexing/chunker"
	"github.com/akshayvibe/GoRag/internal/indexing/embedding"
	"github.com/akshayvibe/GoRag/internal/indexing/upload"
	"github.com/pkoukk/tiktoken-go"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
    log.Fatal("Error loading .env:", err)
	}
	Path := "../data/demo.pdf"
	content, err := upload.OpenPdf(Path)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	//Getting parsed content
	fmt.Printf("Parsed Content: %s", string(data))

	//Get encoding model
	Encoding, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		panic(err)
	}
	//Initialize token chunker
	tokenChunker := chunker.TokenChunker{
		Encoding:    Encoding,
		ChunkSize:   500,
		OverlapSize: 50,
	}
	//Chunk content
	chunks, err := tokenChunker.ChunkText(content.Text, content.Metadata)
	if err != nil {
		panic(err)
	}
	// creating Vector Embeddings for chunks
	texts := make([]string, len(chunks))

	for i := range chunks {
		texts[i] = chunks[i].Content
	}

	vectors, err := embedding.VectorEmbedding(texts)
	if err != nil {
		panic(err)
	}

	for i := range chunks {
		chunks[i].Vector = vectors[i]
	}
	data, err = json.MarshalIndent(chunks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Chunks: %s", string(data))
}
