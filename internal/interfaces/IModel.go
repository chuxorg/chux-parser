package interfaces

// An Interface for Models that interact with a data store
type IModel interface {
	// If the Model has changes, will return true
	IsDirty() bool
	// When the Model is first created,
	// the model is considered New. After the model is
	// Saved or Loaded it is no longer New
	IsNew() bool
	// Saves the Model to a Data Store
	Save() error
	// Loads a Model from the Data Store
	Load(id string) (interface{}, error)
	// Searches for items in the data store
	Search(args ...interface{}) ([]interface{}, error)
	// Deletes a Model from the Data Store
	Delete() error
	// Sets the internal state of the model.
	setState(data []byte)
}
