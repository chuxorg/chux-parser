package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chuxorg/chux-parser/internal/errors"
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

func (f *File) Save(files []interface{}) error {
	database := os.Getenv("MONGO_DATABASE")
	collectionName := "files"
	uri := os.Getenv("MONGO_URI") + os.Getenv("MONGO_DATABASE") + "?retryWrites=true&w=majority"

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return errors.NewChuxParserError("File.Save() Error connecting to MongoDB", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(database).Collection(collectionName)

	res, err := collection.InsertMany(ctx, files)
	if err != nil {
		return errors.NewChuxParserError("File.Save() Error calling InsertMany to MongoDB", err)
	}

	fmt.Printf("Inserted %d documents: %v\n", len(res.InsertedIDs), res.InsertedIDs)
	return nil
}
