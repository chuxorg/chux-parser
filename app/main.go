package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chuxorg/chux-parser/config"
	"github.com/chuxorg/chux-parser/internal/s3"
)

func main() {
	cfg, err := config.LoadConfig("development")
	if err != nil {
		panic(err)
	}
	bucket := s3.New(
		s3.WithConfig(*cfg),
	)

	bucket.DownloadAll()

	// parser := parsing.New(parsing.WithConfig(*cfg))
	// files := getFiles(*cfg)
	// for _, f := range files {
	// 	//"items_sweetwater.com-2023-04-06T21_06_17.291000.jl"
	// 	parser.Parse(f)
	// }
}

func getFiles(cfg config.ParserConfig) []string {
	retVal := []string{}
	dir := cfg.AWS.DownloadPath
	// Walk the directory recursively and search for files with .jl extension
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Check if file extension is .jl
		if filepath.Ext(path) == ".jl" {
			retVal = append(retVal, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return retVal
}
