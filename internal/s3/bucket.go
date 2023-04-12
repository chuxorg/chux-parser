package s3

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/chuxorg/chux-parser/config"
)

const basePath = "data/"

// Define a struct to hold the JSON object's URL field
type Line struct {
	URL string `json:"url"`
}

type IBucket interface {
	getObjects() (*s3.ListObjectsV2Output, error)
	logError(msg string, args ...interface{})
	DownloadAll() []File
	Download(fileName string)
}

type Bucket struct {
	Name         string
	Profile      string
	DownloadPath string
	Session      *session.Session
}

var _cfg *config.ParserConfig

func New(options ...func(*Bucket)) *Bucket {

	bucket := &Bucket{}
	for _, option := range options {
		option(bucket)
	}

	if _cfg != nil {
		bucket.Name = _cfg.AWS.BucketName
		bucket.Profile = _cfg.AWS.Profile
		bucket.DownloadPath = _cfg.AWS.DownloadPath
	}

	return bucket
}

func WithConfig(config config.ParserConfig) func(*Bucket) {
	return func(product *Bucket) {
		_cfg = &config
	}
}

func newFile() *File {
	return &File{}
}

// Logs an Error Message
// Inputs:
//
//	msg is the error message that occurred.
//
// Output:
//
//	None
func (b Bucket) logError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func (b *Bucket) Download() ([]File, error) {
	// Replace with your bucket and region
	s3Bucket := _cfg.AWS.BucketName
	region := os.Getenv("AWS_REGION")

	// Create a new AWS session with the specified region
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	// Create a new S3 service instance using the session
	svc := s3.New(sess)

	// List objects in the S3 bucket
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(s3Bucket)})
	if err != nil {
		fmt.Println("Error listing objects:", err)
		return nil, err
	}
	var files []File
	// Iterate through each object in the bucket
	for _, item := range resp.Contents {
		objectKey := *item.Key

		// Download the object from S3
		fileReader, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    item.Key,
		})
		if err != nil {
			fmt.Println("Error getting object:", err)
			continue
		}

		// Read the first line (JSON object) from the file
		lineReader := bufio.NewReader(fileReader.Body)
		lineStr, err := lineReader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Error reading line:", err)
			continue
		}

		// Unmarshal the JSON object into a Line struct
		var lineObj Line
		err = json.Unmarshal([]byte(lineStr), &lineObj)
		if err != nil {
			fmt.Println("Error un-marshalling JSON:", err)
			continue
		}

		// Extract the FQDN from the URL
		companyName, err := b.extractCompanyName(lineObj.URL)
		if err != nil {
			fmt.Println("Error extracting company name:", err)
			continue
		}

		// Create the company directory if it doesn't exist
		companyPath := filepath.Join(basePath, companyName)
		err = os.MkdirAll(companyPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating company directory:", err)
			continue
		}

		// Download the file to the newly created path
		dateTime, err := b.formatDateDir(objectKey)
		if err != nil {
			fmt.Println("Error creating company directory:", err)
		}
		outputPath := filepath.Join(companyPath, companyName+"-"+dateTime+".jl")
		err = b.downloadFile(fileReader.Body, outputPath)
		if err != nil {
			fmt.Println("Error downloading file:", err)
		}
		// Set a new File Struct to be used during parsing
		file := File{
			Path:         outputPath,
			LastModified: *item.LastModified,
			Size:         *item.Size,
			IsProduct:    true,
			IsParsed:     false,
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		}

		files = append(files, file)
	}
	return files, nil
}

// The extractCompanyName function takes a raw URL string as input, parses it, and extracts the hostname.
// It then uses a regular expression to remove the domain extension and any subdomains (e.g., "www").
// The resulting company name is returned.
func (b *Bucket) extractCompanyName(rawURL string) (string, error) {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	// Extract the hostname from the parsed URL
	host := parsedURL.Hostname()

	// Use a regular expression to match and remove the domain extension
	// This regular expression matches the last segment of the domain name
	// (e.g., "com", "edu") and any subdomains (e.g., "www")
	re := regexp.MustCompile(`(?:\w+\.)?(\w+)\.\w+`)
	matches := re.FindStringSubmatch(host)

	// Check if the regular expression found a match
	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract company name from URL")
	}

	// Return the matched company name (the second element in the matches slice)
	return matches[1], nil
}

// formatDateDir takes an S3 object key string and extracts the date and time
// from the key, formatting it as a directory name like "20230410-20.38.35"
func (b *Bucket) formatDateDir(objectKey string) (string, error) {
	// Use a regular expression to match the date and time from the object key
	// This regular expression looks for a date and time pattern like "2023-04-10T20-38-35"
	re := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})T(\d{2})-(\d{2})-(\d{2})`)
	matches := re.FindStringSubmatch(objectKey)

	// Check if the regular expression found a match
	if len(matches) < 7 {
		return "", fmt.Errorf("could not extract date and time from object key")
	}

	// Extract the matched date and time components
	year := matches[1]
	month := matches[2]
	day := matches[3]
	hour := matches[4]
	minute := matches[5]
	second := matches[6]

	// Format the date and time components as part of the filename name like "20230410T20-38-35"
	dateTimeDir := fmt.Sprintf("%s%s%sT%s-%s-%s", year, month, day, hour, minute, second)

	return dateTimeDir, nil
}

// downloadFile downloads the file from the S3 bucket, reads it from fileReader, and saves it to the specified outputPath.
// Input:
// - fileReader: An io.ReadCloser from which the file content will be read.
// - outputPath: The path where the downloaded file will be saved.
// Returns: An error if something goes wrong during the file download and save process.
func (b *Bucket) downloadFile(fileReader io.ReadCloser, outputPath string) error {
	// Ensure the fileReader is closed after the function finishes.
	defer fileReader.Close()

	// Create a new file or open an existing one at the outputPath.
	// This will also truncate the file if it already exists, meaning
	// it will clear the existing content and set the file size to zero.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	// Ensure the outputFile is closed after the function finishes.
	defer outputFile.Close()

	// Copy the content from fileReader to outputFile.
	// This will write the downloaded data from the S3 bucket into the file.
	_, err = io.Copy(outputFile, fileReader)
	if err != nil {
		return fmt.Errorf("failed to copy data to output file: %w", err)
	}

	return nil
}
