package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/akshayvibe/GoRag/internal/config"
	"github.com/akshayvibe/GoRag/internal/db/vectorstore"
	"github.com/akshayvibe/GoRag/internal/handler"
	"github.com/akshayvibe/GoRag/internal/indexing/chunker"
	"github.com/akshayvibe/GoRag/internal/retrieval"

	"github.com/joho/godotenv"
	"github.com/pkoukk/tiktoken-go"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.Load();
	log.Printf(
		"Qdrant config: host=%q port=%d tls=%t apiKeySet=%t",
		cfg.QdrantHost,
		cfg.QdrantPort,
		cfg.QdrantUseTLS,
		cfg.QdrantAPIKey != "",
	)
	encoding, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		log.Fatalf("Failed to get encoding: %v", err)
	}

	tokenChunker := chunker.TokenChunker{
		Encoding:    encoding,
		ChunkSize:   500,
		OverlapSize: 50,
	}

	store, err := vectorstore.NewQdrantStore(cfg,"demo_collection", 1536)
	if err != nil {
		log.Fatalf("Failed to connect to Qdrant: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	if err := store.CreateCollection(ctx); err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	retriever := retrieval.NewRetriever(store)
	
	appHandler := &handler.AppHandler{
		Store:     store,
		Chunker:   tokenChunker,
		Retriever: retriever,
	}
	http.HandleFunc("/upload", appHandler.UploadEndpoint)
	http.HandleFunc("/query", appHandler.SendReq)
	

	port := ":8080"
	fmt.Printf("Server is starting on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}