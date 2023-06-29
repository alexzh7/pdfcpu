package main

import (
	"log"
	"path/filepath"

	"github.com/alexzh7/pdfcpu/pkg/api"
	"github.com/alexzh7/pdfcpu/pkg/pdfcpu/model"
)

func main() {
	dir := "cmd/example_merge/"
	newFile := filepath.Join(dir, "mergedPdf.pdf")
	err := api.MergeCreateFile([]string{filepath.Join(dir, "111.pdf"), filepath.Join(dir, "222.pdf")}, newFile, model.NewDefaultConfiguration())
	if err != nil {
		panic(err)
	}

	log.Default().Printf("created file %v", newFile)
}
