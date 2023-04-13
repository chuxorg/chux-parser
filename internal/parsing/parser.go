package parsing

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	bo "github.com/chuxorg/chux-models/config"
	"github.com/chuxorg/chux-models/models/articles"
	"github.com/chuxorg/chux-models/models/products"
	"github.com/chuxorg/chux-parser/config"
	cfg "github.com/chuxorg/chux-parser/config"
	"github.com/chuxorg/chux-parser/internal/s3"
)

// Parser struct for parsing
type Parser struct {
	products []products.Product
	articles []articles.Article
}

var _cfg *cfg.ParserConfig
var _bizObjConfig *bo.BizObjConfig

// New returns a new Parser struct
func New(options ...func(*Parser)) *Parser {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	var err error
	_cfg, err = cfg.LoadConfig(env)
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

func WithConfig(config cfg.ParserConfig) func(*Parser) {
	return func(product *Parser) {
		_cfg = &config
		_bizObjConfig = cfg.NewBizObjConfig(_cfg)
	}
}

func (p *Parser) Parse(file s3.File) {

	// Create the out and errOut channels
	out := make(chan string)
	errOut := make(chan error)

	// Call the readJSONObjects function in a separate goroutine
	go readJSONObjects(file.Path, out, errOut)

	// Loop until both channels are closed and set to nil
	for {
		select {
		case jsonStr, ok := <-out:
			if !ok {
				out = nil // Set the channel to nil to stop checking it
			} else {
				// Process the JSON string (e.g., pass it to Product.SetState())
				fmt.Println("JSON Object:", jsonStr)

				if file.IsProduct {
					product := products.New(
						products.WithBizObjConfig(*_bizObjConfig),
					)
					var err error
					err = product.Parse(jsonStr)
					if err != nil {
						fmt.Println(err)
					}
					err = product.Save()
					if err != nil {
						fmt.Println(err)
					}
				} else {
					article := articles.New(
						articles.WithBizObjConfig(*_bizObjConfig),
					)
					err := article.Parse(jsonStr)
					if err != nil {
						fmt.Println(err)
					}
					err = article.Save()
					if err != nil {
						fmt.Println(err)
					}

				}
			}
		case err, ok := <-errOut:
			if !ok {
				errOut = nil // Set the channel to nil to stop checking it
			} else {
				// Handle the error (e.g., log it, exit the program, or take other appropriate action)
				fmt.Println("Error:", err)
			}
		}

		// Break the loop when both channels are closed and set to nil
		if out == nil && errOut == nil {
			break
		}
	}

}

func readJSONObjects(filePath string, out chan<- string, errOut chan<- error) {
	// Close both output channels after the function exits
	defer close(out)
	defer close(errOut)
	megabytes := 5
	byteSize := megabytes * 1024 * 1024
	// Open the specified JSON file
	file, err := os.Open(filePath)
	if err != nil {
		// If an error occurs, send the error to the error output channel
		errOut <- fmt.Errorf("failed to open file: %w", err)
		return
	}
	// Ensure the file is closed after the function exits
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, bufio.MaxScanTokenSize)
	scanner.Buffer(buf, byteSize)
	// Declare a variable to store each JSON object
	var jsonObj map[string]interface{}

	// Skip the first line
	if scanner.Scan() {
	} // Do nothing, just skip the first line

	// Iterate over each line in the file
	for scanner.Scan() {
		// Get the current line as a string
		line := scanner.Text()

		// Unmarshal the JSON line into the jsonObj variable
		err = json.Unmarshal([]byte(line), &jsonObj)
		if err != nil {
			// If an error occurs, send the error to the error output channel
			errOut <- fmt.Errorf("failed to unmarshal JSON object: %w", err)
			continue
		}

		// Marshal the JSON object back to a JSON string
		jsonStr, err := json.Marshal(jsonObj)
		if err != nil {
			// If an error occurs, send the error to the error output channel
			errOut <- fmt.Errorf("failed to marshal JSON object: %w", err)
			continue
		}

		// Send the JSON string to the output channel
		out <- string(jsonStr)
	}

	// Check for any errors that occurred during the scanning process
	if err := scanner.Err(); err != nil {
		// If an error occurs, send the error to the error output channel
		errOut <- fmt.Errorf("error scanning file: %w", err)
	}
}

func (p *Parser) GetFiles(cfg config.ParserConfig) []string {
	retVal := []string{}
	dir := cfg.AWS.DownloadPath
	// Walk the directory recursively and search for files with .jl extension
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Check if file extension is .jl
		if filepath.Ext(path) == ".jl" {
			retVal = append(retVal, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return retVal
}
