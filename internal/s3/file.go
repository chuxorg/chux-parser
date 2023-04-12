package s3

import (
	"encoding/json"
	"time"
)

// The file Struct is used to track the status of a file's
// parsing and storage in the datastores
type File struct {
	Path         string    `json:"path"`
	LastModified time.Time `json:"lastModified"`
	Size         int64     `json:"size"`
	IsParsed     bool      `json:"isParsed"`
	IsProduct    bool      `json:"isProduct"`
	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
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
