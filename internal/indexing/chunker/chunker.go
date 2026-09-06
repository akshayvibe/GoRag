package chunker

import (
	"github.com/pkoukk/tiktoken-go"
	"fmt"
	"github.com/akshayvibe/GoRag/internal/models"
)


type TokenChunker struct{
	Encoding  *tiktoken.Tiktoken
	ChunkSize int
	OverlapSize int
}
func NewTokenChunker(encodingName string, chunkSize int, overlapSize int) *TokenChunker {	
	tke, err := tiktoken.GetEncoding(encodingName)
	if chunkSize<=0 || overlapSize<0 {
		return nil
	}
	if tke == nil {
		return nil
	}
	if overlapSize >= chunkSize {
		return nil
	}
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return nil
	}
	return &TokenChunker{
		Encoding:  tke,
		ChunkSize: chunkSize,
		OverlapSize: overlapSize,
	}
}

func (tc *TokenChunker) ChunkText(text string,metadata models.Metadata) ([]*models.Chunk,error) {	
	if tc == nil || tc.Encoding == nil {
			return nil, fmt.Errorf("chunker is not initialized")
		}
	
		tokens := tc.Encoding.Encode(text, nil, nil)
	
		if len(tokens) == 0 {
			return []*models.Chunk{}, nil
		}
	
		step := tc.ChunkSize - tc.OverlapSize
	
		chunks := make([]*models.Chunk, 0)
	
		for start := 0; start < len(tokens); start += step {
			end := start + tc.ChunkSize
	
			if end > len(tokens) {
				end = len(tokens)
			}
	
			chunkTokens := tokens[start:end]
	
			chunks = append(chunks, &models.Chunk{
				ID:       fmt.Sprintf("chunk-%d", len(chunks)),
				Content:  tc.Encoding.Decode(chunkTokens),
				Tokens:   len(chunkTokens),
				Metadata: metadata,
			})
	
			if end == len(tokens) {
				break
			}
		}
	
		return chunks, nil
}