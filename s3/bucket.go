package s3

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/chuxorg/chux-parser/errors"
	"github.com/chuxorg/chux-parser/logging"
	"golang.org/x/net/publicsuffix"
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

func New(options ...func(*Bucket)) *Bucket {
	logging.Debug("Creating new Bucket struct")
	bucket := &Bucket{}
	for _, option := range options {
		option(bucket)
	}
	bucketName := os.Getenv("AWS_SOURCE_BUCKET")
	bucket.Name = bucketName
	bucket.Profile = "csailer"
	bucket.DownloadPath = os.Getenv("AWS_DOWNLOAD_PATH")
	logging.Debug("Bucket struct created with the following settings\nName: %s\nProfile: %s\nDownloadPath: %s", bucket.Name, bucket.Profile, bucket.DownloadPath)
	return bucket
}

func (b *Bucket) Download() ([]File, error) {
	logging.Debug("Bucket.Download() called")
	s3Bucket := os.Getenv("AWS_SOURCE_BUCKET")
	region := "us-east-1"
	logging.Info("Downloading files from S3 bucket %s", s3Bucket)
	// Create a new AWS session with the specified region
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	// Create a new S3 service instance using the session
	svc := s3.New(sess)

	// List objects in the S3 bucket
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(s3Bucket)})
	if err != nil {
		logging.Error("Bucket.Download() Error listing objects:", err)
		return nil, errors.NewChuxParserError("Bucket.Download() Error listing objects:", err)
	}
	var files []File

	for _, item := range resp.Contents {
		// Download the object from S3
		fileReader, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    item.Key,
		})
		if err != nil {
			msg := fmt.Sprintf("Bucket.Download() Error getting object %v. Continuing", err)
			logging.Warning(msg, err)
			continue
		}

		lineReader := bufio.NewReader(fileReader.Body)
		lineStr, err := lineReader.ReadString('\n')
		if err != nil && err != io.EOF {
			msg := fmt.Sprintf("Bucket.Download() Error reading line:%v. Continuing", err)
			logging.Warning(msg, err)
			continue
		}

		// Unmarshal the JSON object into a Line struct
		var lineObj Line
		err = json.Unmarshal([]byte(lineStr), &lineObj)
		if err != nil {
			msg := fmt.Sprintf("Bucket.Download() Error unmarshalling JSON object:%v. Continuing", err)
			logging.Warning(msg, err)
			continue
		}

		// Extract the FQDN from the URL
		companyName, err := b.extractCompanyName(lineObj.URL)
		if err != nil {
			msg := fmt.Sprintf("Bucket.Download() Error extracting company name:%v. Continuing", err)
			log.Println(msg, err)
			continue
		}

		// Read the entire content of the file
		contentBytes, err := ioutil.ReadAll(fileReader.Body)
		if err != nil {
			msg := fmt.Sprintf("Bucket.Download() Error reading file content:%v. Continuing", err)
			log.Println(msg, err)
			continue
		}

		content := string(contentBytes)
		if !strings.Contains(strings.ToLower(companyName), "ebay") {
			file := File{
				Content:      content,
				LastModified: *item.LastModified,
				Size:         *item.Size,
				IsProduct:    b.isProduct(companyName),
				IsParsed:     false,
				DateCreated:  time.Now(),
				DateModified: time.Now(),
			}

			files = append(files, file)
		}
	}
	logging.Info("Bucket.Download() Files Ready to Process: ", len(files))
	return files, nil
}

// The extractCompanyName function takes a raw URL string as input, parses it, and extracts the hostname.
// It then removes the domain extension and any subdomains (e.g., "www").
// The resulting company name is returned.
func (b *Bucket) extractCompanyName(rawURL string) (string, error) {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	// Extract the hostname from the parsed URL
	host := parsedURL.Hostname()

	// Use the public suffix library to obtain the effective TLD plus one
	etldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return "", err
	}

	// Remove the TLD plus one from the hostname
	trimmedHost := strings.TrimSuffix(host, "."+etldPlusOne)

	// Split the remaining hostname by dots and take the last part as the company name
	parts := strings.Split(trimmedHost, ".")
	companyName := parts[len(parts)-1]

	return companyName, nil
}

// The isProduct function takes a slice of strings and a target string as input.
func (b *Bucket) isProduct(target string) bool {

	productSources := []string{
		"ebay",
		"sweetwater",
		"perfectcircuit",
		"reverb",
		"thomann",
		"zzounds",
		"samash",
		"guitarcenter",
		"musiciansfriend",
		"thomannmusic",
	}

	target = strings.ToLower(strings.TrimSpace(target))

	for _, value := range productSources {
		if value == target {
			return true
		}
	}
	return false
}
