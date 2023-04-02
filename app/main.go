package main

import (
	"github.com/chuxorg/chux-parser/config"
	"github.com/chuxorg/chux-parser/internal/parsing"
)

func main() {
	cfg, err := config.LoadConfig("development")
	if err != nil {
		panic(err)
	}
	parser := parsing.New(parsing.WithConfig(*cfg))
	parser.Parse("items_guitarcenter.com-2023-03-20T18_53_00.897000.json")
}
