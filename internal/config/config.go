package config

import (
	"os"
	"strconv"
)

type Config struct {
	QdrantHost   string
	QdrantPort   int
	QdrantAPIKey string
	QdrantUseTLS bool
	// OllamaURL   string
}

func Load() Config {
	// ollamaURL := os.Getenv("OLLAMA_URL")
	port, _ := strconv.Atoi(os.Getenv("QDRANT_PORT"))
	useTLS, _ := strconv.ParseBool(os.Getenv("QDRANT_USE_TLS"))

	return Config{
		QdrantHost:   os.Getenv("QDRANT_HOST"),
		QdrantPort:   port,
		QdrantAPIKey: os.Getenv("QDRANT_API_KEY"),
		QdrantUseTLS: useTLS,
		// OllamaURL:   ollamaURL,
	}
}