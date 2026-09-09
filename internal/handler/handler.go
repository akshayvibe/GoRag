package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/akshayvibe/GoRag/internal/db/vectorstore"
	"github.com/akshayvibe/GoRag/internal/indexing/chunker"
	"github.com/akshayvibe/GoRag/internal/indexing/embedding"
	"github.com/akshayvibe/GoRag/internal/indexing/upload"
	"github.com/akshayvibe/GoRag/internal/models"
)

type AppHandler struct {
	Store   *vectorstore.QdrantStore
	Chunker chunker.TokenChunker
}

// UploadEndpoint is now a clean orchestrator.
func (h *AppHandler) UploadEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Extract and save the file temporarily
	tempPath, originalName, err := h.saveUploadedFile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer os.Remove(tempPath) // Clean up automatically when done

	// 2. Parse the PDF and chunk the text
	chunks, err := h.parseAndChunk(tempPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Generate embeddings and store them in Qdrant
	if err := h.embedAndStore(r.Context(), chunks, originalName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return success response
	h.respondSuccess(w, originalName, len(chunks))
}

// --- Helper Methods ---

// saveUploadedFile handles form parsing and writes the upload to a temp file.
func (h *AppHandler) saveUploadedFile(r *http.Request) (string, string, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return "", "", fmt.Errorf("failed to parse form: %w", err)
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", "", fmt.Errorf("failed to retrieve file: %w", err)
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, file); err != nil {
		return "", "", fmt.Errorf("failed to save file: %w", err)
	}

	return tempFile.Name(), fileHeader.Filename, nil
}

// parseAndChunk extracts text from the PDF and cuts it into chunks.
func (h *AppHandler) parseAndChunk(filePath string) ([]*models.Chunk, error) {
	pdfContent, err := upload.OpenPdf(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PDF: %w", err)
	}

	chunks, err := h.Chunker.ChunkText(pdfContent.Text, pdfContent.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to chunk text: %w", err)
	}

	return chunks, nil
}

// embedAndStore fetches vector embeddings and saves them to the database.
func (h *AppHandler) embedAndStore(ctx context.Context, chunks []*models.Chunk, filename string) error {
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}

	vectors, err := embedding.VectorEmbedding(ctx, texts)
	if err != nil {
		return fmt.Errorf("failed to create embeddings: %w", err)
	}

	ids := make([]uint64, len(chunks))
	insertSources := make([]string, len(chunks))

	for i := range chunks {
		ids[i] = uint64(i + 1)
		insertSources[i] = filename
	}

	if err := h.Store.AddDocuments(ctx, ids, vectors, texts, insertSources); err != nil {
		return fmt.Errorf("failed to store in Qdrant: %w", err)
	}

	return nil
}

// respondSuccess formats the JSON response.
func (h *AppHandler) respondSuccess(w http.ResponseWriter, filename string, chunkCount int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "File successfully indexed!",
		"chunks":  chunkCount,
		"file":    filename,
	})
}