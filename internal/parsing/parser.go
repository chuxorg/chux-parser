package parsing

import (
	"log"
	"os"

	"github.com/chuxorg/chux-models/models/articles"
	"github.com/chuxorg/chux-models/models/products"
	"github.com/chuxorg/chux-parser/config"
)

// Parser struct for parsing
type Parser struct {
	products []products.Product
	articles []articles.Article
}

var _cfg *config.ParserConfig

// New returns a new Parser struct
func New(options ...func(*Parser)) *Parser {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	var err error
	_cfg, err = config.LoadConfig(env)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	parser := &Parser{}
	for _, option := range options {
		option(parser)
	}
	return parser
}

func WithLoggingLevel(level string) func(*Parser) {
	return func(product *Parser) {
		_cfg..Level = level
	}
}
