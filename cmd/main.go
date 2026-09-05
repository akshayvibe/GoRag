package main

import (
	"fmt"
	"github.com/akshayvibe/GoRag/internal/ingestion"
)

func main() {
	//getting the pdf content from the terminal
	Path:="../data/demo.pdf"
	content,err:=ingestion.OpenPdf(Path)
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
}
