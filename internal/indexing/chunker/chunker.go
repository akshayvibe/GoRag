package chunker

import (
	"github.com/pkoukk/tiktoken-go"
	"fmt"
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