package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/chuxorg/chux-parser/logging"
	"github.com/chuxorg/chux-parser/parsing"
	"github.com/chuxorg/chux-parser/s3"
)

var logFileMutex sync.Mutex
var logFile *os.File
var logger *logging.Logger

func main() {

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("LOG_LEVEL", "0")
	logger = logging.NewLogger(logging.LogLevelDebug)
	err := fetchAndSetSecrets("dev/secrets")
	if err != nil {
		log.Fatalf("failed to fetch and set secrets: %v", err)
	}
	fmt.Print("Setting up logging...")
	setUpLogging()
	defer closeLogFile()
	logger.Debug("Logging set up")
	logger.Info("Logging set up")
	bucket := s3.New(
		s3.WithLogger(logger),
	)

	files, err := bucket.Download()
	if err != nil {
		logger.Error("Failed to download files from S3", err)
		panic(err)
	}

	parser := parsing.New(
		parsing.WithLogger(logger),
	)
	logger.Info("Parsing %d Products and Articles", len(files))
	startTime := time.Now()
	for _, f := range files {
		parser.Parse(f)
	}
	elapsedTime := time.Since(startTime).Seconds()
	logger.Info("Parsed %d Articles and Products in %d seconds", len(files), elapsedTime)

	filesInterface := make([]interface{}, len(files))
	for i, file := range files {
		filesInterface[i] = file
	}
	file := s3.NewFile(
		s3.FileWithLogger(logger),
	)
	logger.Info("Saving information of %d file to MongoDB", len(files))
	err = file.Save(filesInterface)
	if err != nil {
		logger.Error("Failed to save files to MongoDB", err)
		panic(err)
	}
}

func setUpLogging() {

	var err error

	logDir := "logs/chux-cprs/"
	err = os.MkdirAll(logDir, 0755) // Set permissions to 0755
	if err != nil {
		log.Fatalf("Error creating log directory: %v", err)
	}

	// Open the log file
	logFilePath := filepath.Join(logDir, "chux-parser.log")
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = int(logging.LogLevelInfo)
	}

	logger = logging.NewLogger(logging.LogLevel(logLevel))
	logger.SetOutput(logFile)
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

func closeLogFile() {
	logFileMutex.Lock()
	defer logFileMutex.Unlock()

	if logFile != nil {
		logFile.Close()
	}
}
