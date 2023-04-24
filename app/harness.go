package app

import (
	"github.com/chuxorg/chux-parser/parsing"
	"github.com/chuxorg/chux-parser/s3"
)

func TestHarness() {

	bucket := s3.New()

	files, err := bucket.Download()
	if err != nil {
		panic(err)
	}

	parser := parsing.New()
	for _, f := range files {
		parser.Parse(f)
	}
	filesInterface := make([]interface{}, len(files))
	for i, file := range files {
		filesInterface[i] = file
	}
	file := s3.File{}
	file.Save(filesInterface)
}
