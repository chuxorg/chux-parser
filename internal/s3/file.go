package s3

import (
	"encoding/json"
	"time"
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
