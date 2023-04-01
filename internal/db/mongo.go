package db

func MongoDB() {
	config := mongodb.MongoConfig{
		CollectionName: "testCollection",
		DatabaseName:   "testDatabase",
		URI:            "mongodb://localhost:27017",
	}
	mdb := mongodb.NewMongoDB(config)
	return mdb
}
