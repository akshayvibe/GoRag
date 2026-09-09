package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/akshayvibe/GoRag/internal/db/vectorstore"
	"github.com/akshayvibe/GoRag/internal/handler"
	"github.com/akshayvibe/GoRag/internal/indexing/chunker"
	"github.com/joho/godotenv"
	"github.com/pkoukk/tiktoken-go"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: No .env file found or error loading it.")
	}

	encoding, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		log.Fatalf("Failed to get encoding: %v", err)
	}

	tokenChunker := chunker.TokenChunker{
		Encoding:    encoding,
		ChunkSize:   500,
		OverlapSize: 50,
	}

	store, err := vectorstore.NewQdrantStore("localhost", 6334, "demo_collection", 1536)
	if err != nil {
		log.Fatalf("Failed to connect to Qdrant: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	if err := store.CreateCollection(ctx); err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	appHandler := &handler.AppHandler{
		Store:   store,
		Chunker: tokenChunker,
	}

	http.HandleFunc("/upload", appHandler.UploadEndpoint)

	port := ":8080"
	fmt.Printf("Server is starting on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}