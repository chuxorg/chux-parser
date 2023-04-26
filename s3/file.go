package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/chuxorg/chux-parser/errors"
	"github.com/chuxorg/chux-parser/logging"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// The file Struct is used to track the status of a file's
// parsing and storage in the datastores
type File struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OwnerID      primitive.ObjectID `bson:"ownerId,omitempty" json:"ownerId,omitempty"`
	Company      string             `bson:"company, omitempty" json:"company,omitempty"`
	Content      string             `bson:"content", omitempty" json:"content,omitempty"`
	LastModified time.Time          `bson:"lastModified", omitempty" json:"lastModified,omitempty"`
	Size         int64              `bson:"size", omitempty" json:"size,omitempty"`
	IsProduct    bool               `bson:"isProduct", omitempty" json:"isProduct,omitempty"`
	IsParsed     bool               `bson:"isParsed", omitempty" json:"isParsed,omitempty"`
	DateCreated  time.Time          `bson:"dateCreated", omitempty" json:"dateCreated,omitempty"`
	DateModified time.Time          `bson:"dateModified", omitempty" json:"dateModified,omitempty"`
	Path         string             `bson:"path", omitempty" json:"path,omitempty"`
	ArchivedPath string             `bson:"archivedPath", omitempty" json:"archivedPath,omitempty"`
	Logger       *logging.Logger    `bson:"-" json:"-"`
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
	logging.Info("Inserting %d files to MongoDB", len(files))
	cnt := 0
	for _, file := range files {
		f, ok := file.(File)
		if !ok {
			// handle error when file is not of type File
			continue
		}
		f.Content = ""
		_, err := collection.InsertOne(ctx, f)
		if err != nil {
			logging.Error("File.Save() error calling InsertOne to MongoDB", err)
			continue
		}
		cnt++
	}

	logging.Info("Inserted %d File documents.\n", cnt)

	return nil
}
