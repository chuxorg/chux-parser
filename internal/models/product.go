package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Offer struct {
	Price        string `bson:"price"`
	Currency     string `bson:"currency"`
	Availability string `bson:"availability"`
}

type Breadcrumb struct {
	Name string `bson:"name"`
	Link string `bson:"link"`
}

type AdditionalProperty struct {
	Name  string `bson:"name"`
	Value string `bson:"value"`
}

type AggregateRating struct {
	RatingValue float64 `bson:"ratingValue"`
	BestRating  float64 `bson:"bestRating"`
	ReviewCount int     `bson:"reviewCount"`
}

type Product struct {
	ID                   primitive.ObjectID   `bson:"_id,omitempty"`
	URL                  string               `bson:"url"`
	CanonicalURL         string               `bson:"canonicalUrl"`
	Probability          float64              `bson:"probability"`
	Name                 string               `bson:"name"`
	Offers               []Offer              `bson:"offers"`
	SKU                  string               `bson:"sku"`
	MPN                  string               `bson:"mpn"`
	Brand                string               `bson:"brand"`
	Breadcrumbs          []Breadcrumb         `bson:"breadcrumbs"`
	MainImage            string               `bson:"mainImage"`
	Images               []string             `bson:"images"`
	Description          string               `bson:"description"`
	DescriptionHTML      string               `bson:"descriptionHtml"`
	AdditionalProperties []AdditionalProperty `bson:"additionalProperty"`
	AggregateRating      AggregateRating      `bson:"aggregateRating"`
	Color                string               `bson:"color"`
	Style                string               `bson:"style"`
	DateCreated          CustomTime           `bson:"dateCreated"`
	DateModified         CustomTime           `bson:"dateModified"`
	isNew                bool                 `bson:"isNew"`
	isDeleted            bool                 `bson:"isDeleted"`
	isDirty              bool                 `bson:"isDirty"`
}

func NewProduct() *Product {
	return &Product{}
}

func Create() {
	mdb := MongoDB()

}

func (p *Product) Serialize() ([]byte, error) {
	return Marshal(p)
}

func (p *Product) Deserialize(jsonData []byte) (interface{}, error) {
	var product Product
	return Unmarshal(jsonData, product)
}

func MongoDB() {
	config := mongodb.MongoConfig{
		CollectionName: "testCollection",
		DatabaseName:   "testDatabase",
		URI:            "mongodb://localhost:27017",
	}
	mdb := mongodb.NewMongoDB(config)
	return mdb
}
