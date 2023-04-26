package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/chuxorg/chux-parser/errors"
	"github.com/chuxorg/chux-parser/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// The file Struct is used to track the status of a file's
// parsing and storage in the datastores
type File struct {
	Content      string
	LastModified time.Time
	Size         int64
	IsProduct    bool
	IsParsed     bool
	DateCreated  time.Time
	DateModified time.Time
	Path         string
	ArchivedPath string
	Logger       *logging.Logger
}

func NewFile(options ...func(*File)) *File {

	file := &File{}
	for _, option := range options {
		option(file)
	}
	file.DateCreated = time.Now()
	file.DateModified = time.Now()
	file.IsParsed = false
	file.IsProduct = false
	return file
}

func FileWithLogger(logger *logging.Logger) func(*File) {
	return func(file *File) {
		file.Logger = logger
	}
}

func (f *File) ToString() string {
	return f.Path
}

func (f *File) ToJSON() string {

	data, err := json.Marshal(f)
	if err != nil {
		return ""
	}
	return string(data)
}

// Save saves the file to the MongoDB database using InsertMany (bulk insert)
func (f *File) Save(files []interface{}) error {
	logging := f.Logger
	logging.Debug("File.Save() called")
	database := os.Getenv("MONGO_DATABASE")
	collectionName := "files"
	username := os.Getenv("MONGO_USER_NAME")
	password := os.Getenv("MONGO_PASSWORD")

	uri := os.Getenv("MONGO_URI")
	mongoURL := fmt.Sprintf(uri, username, password)
	uri = mongoURL + os.Getenv("MONGO_DATABASE") + "?retryWrites=true&w=majority"

	maskedURI := fmt.Sprintf(uri, "******", "******") + "?retryWrites=true&w=majority"
	logging.Info("Saving to MongoDB: %s", maskedURI)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		logging.Error("File.Save() error creating new client", err)
		return errors.NewChuxParserError("File.Save() Error creating new client", err)
	}

	logging.Info("Connecting to MongoDB with a 45 second timeout")
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logging.Error("File.Save() error connecting to MongoDB", err)
		return errors.NewChuxParserError("File.Save() Error connecting to MongoDB", err)
	}

	logging.Info("Connected to MongoDB")

	defer client.Disconnect(ctx)

	collection := client.Database(database).Collection(collectionName)
	logging.Info("Bulk Inserting %d files to MongoDB", len(files))
	res, err := collection.InsertMany(ctx, files)
	if err != nil {
		logging.Error("File.Save() error calling InsertMany to MongoDB", err)
		return errors.NewChuxParserError("File.Save() Error calling InsertMany to MongoDB", err)
	}
	logging.Info("Bulk inserted %d File documents.\n", len(res.InsertedIDs))

	return nil
}
