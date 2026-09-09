package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/akshayvibe/GoRag/internal/config"
	"github.com/akshayvibe/GoRag/internal/handler"
	"github.com/akshayvibe/GoRag/internal/indexing/chunker"
	"github.com/akshayvibe/GoRag/internal/indexing/embedding"
	"github.com/akshayvibe/GoRag/internal/indexing/upload"
	"github.com/joho/godotenv"
	"github.com/pkoukk/tiktoken-go"
)

func main() {

	// cfg:=config.Load();
	err := godotenv.Load("../.env")
	if err != nil {
    log.Fatal("Error loading .env:", err)
	}
	Path := "../data/demo.pdf"
	PdfData, err := handler.UploadHandler(Path)
	if err != nil {
		panic(err)
	}
	//Getting parsed content
	fmt.Printf("Parsed Content: %s", PdfData)

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

	vectors, err := embedding.VectorEmbedding(context.Background(),texts)
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
