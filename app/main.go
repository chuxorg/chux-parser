package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chuxorg/chux-parser/config"
	"github.com/chuxorg/chux-parser/internal/s3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := config.LoadConfig("development")
	if err != nil {
		panic(err)
	}

    cfg.AWS.AccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	cfg.AWS.SecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

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
