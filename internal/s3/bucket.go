package s3

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/chuxorg/chux-parser/config"
)

type File struct {
	Name         string
	LastModified time.Time
	Size         int64
	IsProduct    bool
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

// GetObjects retrieves the objects in an Amazon S3 bucket
// Inputs:
//
//	sess is the current session, which provides configuration for the SDK's service clients
//
// Output:
//
//	If success, the list of objects and nil
//	Otherwise, nil and an error from the call to ListObjectsV2
func (b *Bucket) getObjects() (*s3.ListObjectsV2Output, error) {

	svc := s3.New(b.Session)
	if b.Name == "" {
		b.logError("GetObjects: Bucket Name is not set")
	}
	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: &b.Name})

	if err != nil {
		return nil, err
	}

	return resp, nil
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

// Starts an AWS Session if one does not
// exist and sets the Session field of
// the Bucket struct
// Inputs:
//
//	None
//
// Output:
//
//	None
func (b *Bucket) startSession() {
	if b.Session != nil {
		// Session already started
		return
	}
	// create the session
	var err error
	// sess, err := session.NewSessionWithOptions(session.Options{
	// })
	sess, err := session.NewSession()
	
	if err != nil {
		b.logError("GetSession: couldn't create an AWS session", err)
	}
	// assign the session to the Bucket
	b.Session = sess
}

// Downloads all items in an s3 bucket.
// The name of the Bucket is assigned to the Bucket Object
// Returns an Slice of file names
// Inputs:
//
//	None
//
// Output:
//
//	A slice of File structs that describes the items that were downloaded
func (b *Bucket) DownloadAll() []File {

	if b.Session == nil {
		b.startSession()
	}

	retVal := []File{}

	resp, err := b.getObjects()
	if err != nil {
		b.logError("DownloadAll: Error from getObjects", err)
	}
	// go over the contents ...
	for _, item := range resp.Contents {
		// create a new File struct
		file := newFile()
		// set the fields of file
		file.Name = *item.Key
		file.LastModified = *item.LastModified
		file.Size = *item.Size
		file.IsProduct = !strings.Contains(*item.Key, "article")
		// append to the slice of Files
		retVal = append(retVal, *file)
		// Download the file
		// if the downloadpath doesn't exist, create the directory
		if _, err := os.Stat(b.DownloadPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(b.DownloadPath, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
		fileName := b.DownloadPath + "/" + file.Name
		b.Download(&fileName)
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", b.Name)
	// return a slice of files that were downloaded
	return retVal
}

// Downloads an Object from an s3 Bucket
// Inputs:
//
//	objectName is the name of the Object to download from the Bucket
//
// Output:
//
//	Returns an error if one occurs. Otherwise, returns nil
func (b *Bucket) Download(objectName *string) error {
	fmt.Println("Downloading ", *objectName)
	fileName := path.Base(*objectName)
	file, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(b.Session)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: &b.Name,
			Key:    objectName,
		})

	if err != nil {
		return err
	}

	return nil
}
