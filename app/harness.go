package app

import (
	"log"
	"os"

	"github.com/chuxorg/chux-parser/config"
	"github.com/chuxorg/chux-parser/internal/parsing"
	"github.com/chuxorg/chux-parser/internal/s3"
	"github.com/joho/godotenv"
)

func TestHarness() {
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

	bucket := s3.New()

	files, err := bucket.Download()
	if err != nil {
		panic(err)
	}

	parser := parsing.New()
	for _, f := range files {
		parser.Parse(f)
	}
	filesInterface := make([]interface{}, len(files))
	for i, file := range files {
		filesInterface[i] = file
	}
	file := s3.File{}
	file.Save(filesInterface)
}
