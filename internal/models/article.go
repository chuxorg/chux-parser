package models

import (
	"github.com/csailer/chux-mongo/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create the MongoDB from chux-mongo
var mongoDB = &db.MongoDB{}

// The Article struct represents an Article Document in MongoDB
type Article struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	URL              string             `bson:"url"`
	Probability      float64            `bson:"probability"`
	Headline         string             `bson:"headline"`
	DatePublished    CustomTime         `bson:"datePublished"`
	DatePublishedRaw string             `bson:"datePublishedRaw"`
	DateCreated      CustomTime         `bson:"dateCreated"`
	DateModified     CustomTime         `bson:"dateModified"`
	DateModifiedRaw  string             `bson:"dateModifiedRaw"`
	Author           string             `bson:"author"`
	AuthorsList      []string           `bson:"authorsList"`
	InLanguage       string             `bson:"inLanguage"`
	Breadcrumbs      []Breadcrumb       `bson:"breadcrumbs"`
	MainImage        string             `bson:"mainImage"`
	Images           []string           `bson:"images"`
	Description      string             `bson:"description"`
	ArticleBody      string             `bson:"articleBody"`
	ArticleBodyHTML  string             `bson:"articleBodyHtml"`
	CanonicalURL     string             `bson:"canonicalUrl"`
	isNew            bool               `bson:"isNew"`
	isDeleted        bool               `bson:"isDeleted"`
	isDirty          bool               `bson:"isDirty"`
}

func NewArticle() *Article {

	return &Article{
		isDirty:   false,
		isNew:     true,
		isDeleted: false,
	}
}

func (a *Article) GetCollectionName() string {
	return "articles"
}

func (a *Article) GetDatabaseName() string {
	return "chux-mongo"
}

func (a *Article) GetURI() string {
	return "mongodb://localhost:27017"
}

// If the Model has changes, will return true
func (a *Article) IsDirty() bool {
	return a.isDirty
}

// When the Model is first created,
// the model is considered New. After the model is
// Saved or Loaded it is no longer New
func (a *Article) IsNew() bool {
	return a.isNew
}

// Saves the Model to a Data Store
func (a *Article) Save() error {
	if a.isNew {
		//--Create a new document
		err := mongoDB.Create(a)
		if err != nil {
			return err
		}

	} else if a.isDirty && !a.isDeleted {
		//--update this document
		err := mongoDB.Update(a, a.ID.String())
		if err != nil {
			return err
		}
	} else if a.isDeleted && !a.isNew {
		//--delete the document
		err := mongoDB.Delete(a, a.ID.String())
		if err != nil {
			return err
		}
	}
	// reset flags
	a.isNew = false
	a.isDirty = false
	a.isDeleted = false
	return nil
}

// Loads a Model from the Data Store
func (a *Article) Load(id string) (interface{}, error) {
	retVal, err := mongoDB.GetByID(a, id)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}

// Deletes a Model from the Data Store
func (a *Article) Delete() error {
	a.isDeleted = true
	return nil
}

// Sets the internal state of the model.
func (a *Article) setState(data []byte) {

}

func (a *Article) Search(args ...interface{}) ([]interface{}, error) {
	return nil, nil
}

func (a *Article) Serialize() ([]byte, error) {
	return Marshal(a)
}

func (a *Article) Deserialize(jsonData []byte) (interface{}, error) {
	var article Article
	return Unmarshal(jsonData, &article)
}
