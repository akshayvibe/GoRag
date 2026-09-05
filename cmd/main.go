package main

import (
	"fmt"


	"github.com/akshayvibe/GoRag/internal/ingestion/upload"
)

func main() {

	// Part -1 ingestion pipeline
	Path:="../data/demo.pdf"
	content,err:=upload.OpenPdf(Path)
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
}

