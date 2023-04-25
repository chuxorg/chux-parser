package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	dsl "github.com/chuxorg/chux-datastore/logging"
	ml "github.com/chuxorg/chux-models/logging"
	"github.com/chuxorg/chux-parser/logging"
	pl "github.com/chuxorg/chux-parser/logging"
	"github.com/chuxorg/chux-parser/parsing"
	"github.com/chuxorg/chux-parser/s3"
)

func main() {

	err := fetchAndSetSecrets("dev/secrets")
	if err != nil {
		log.Fatalf("failed to fetch and set secrets: %v", err)
	}
	fmt.Print("Setting up logging...")
	setUpLogging()

	dsl.Info("Logging set up")
	bucket := s3.New()

	files, err := bucket.Download()
	if err != nil {
		logging.Error("Failed to download files from S3", err)
		panic(err)
	}

	parser := parsing.New()
	dsl.Info("Parsing %d Products and Articles", len(files))
	startTime := time.Now()
	for _, f := range files {
		parser.Parse(f)
	}
	elapsedTime := time.Since(startTime).Seconds()
	dsl.Info("Parsed %d Articles and Products in %d seconds", len(files), elapsedTime)

	filesInterface := make([]interface{}, len(files))
	for i, file := range files {
		filesInterface[i] = file
	}
	file := s3.File{}
	dsl.Info("Saving information of %d file to MongoDB", len(files))
	err = file.Save(filesInterface)
	if err != nil {
		dsl.Error("Failed to save files to MongoDB", err)
		panic(err)
	}
}

func setUpLogging() {

	logFile, err := os.OpenFile("chux-parser.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	pl.SetOutput(logFile)
	ml.SetOutput(logFile)
	dsl.SetOutput(logFile)

	logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = int(pl.LogLevelInfo)
	}

	pl.SetLogLevel(pl.LogLevel(logLevel))
	ml.SetLogLevel(ml.LogLevel(logLevel))
	dsl.SetLogLevel(dsl.LogLevel(logLevel))
}

func fetchAndSetSecrets(secretID string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return fmt.Errorf("failed to get secret value: %v", err)
	}

	var secrets map[string]string
	err = json.Unmarshal([]byte(*result.SecretString), &secrets)
	if err != nil {
		return fmt.Errorf("failed to unmarshal secrets: %v", err)
	}

	for key, value := range secrets {
		err = os.Setenv(key, value)
		fmt.Printf("Setting env variable: %s\n", key)
		if err != nil {
			return fmt.Errorf("failed to set environment variable: %v", err)
		}
	}

	return nil
}
