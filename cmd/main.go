package main

import (
	"fmt"
	"log"
	"encoding/json"

	"github.com/akshayvibe/GoRag/internal/indexing/upload"
)

func main() {
	// Part -1 indexing pipeline
	Path:="../data/demo.pdf"
	content,err:=upload.OpenPdf(Path)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(content, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
}

