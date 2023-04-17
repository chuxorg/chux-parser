package main

import (
	"log"
	"os"

	"github.com/chuxorg/chux-parser/config"
	"github.com/chuxorg/chux-parser/internal/parsing"
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

	/* 
	 Download all files from S3 bucket
	 Parse all files. This will create a new Product or Article for each line in the downloaded files
	 Update the categories for each Product or Article. This will create a new Category for each unique category in the Product or Article
	 Update/Create Company names in MongoDB
	 Message the Image Service to download and process images via Kafka topic
	 Messaage the Cache Service to update the cache via Kafka topic
	 Move processed files to the processed s3 bucket
	 Delete files from the downloaded s3 bucket
	*/

	// bucket := s3.New(
	//  	s3.WithConfig(*cfg),
	// )

	// files, err := bucket.Download()
	// if err != nil {
	// 	panic(err)
	// }

	parser := parsing.New(parsing.WithConfig(*cfg))
	
	// for _, f := range files {
	// 	parser.Parse(f)
	// }
	err = parser.UpdateCategories()
}
