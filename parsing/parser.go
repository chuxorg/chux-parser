package parsing

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	ml "github.com/chuxorg/chux-models/logging"
	"github.com/chuxorg/chux-models/models"
	"github.com/chuxorg/chux-parser/logging"
	"github.com/chuxorg/chux-parser/s3"
)

// Parser struct for parsing
type Parser struct {
	products []models.Product
	articles []models.Article
	Logger   *logging.Logger
}

// New returns a new Parser struct
func New(options ...func(*Parser)) *Parser {

	parser := &Parser{}
	for _, option := range options {
		option(parser)
	}
	parser.Logger.Debug("Creating new Parser struct")
	return parser
}

func WithLogger(logger *logging.Logger) func(*Parser) {
	return func(parser *Parser) {
		parser.Logger = logger
	}
}

func (p *Parser) Parse(file s3.File) {

	productCount := 0
	articleCount := 0
	modelsLogger := ml.NewLogger(ml.LogLevelDebug)
	p.Logger.Debug("Parser.Parse() called")
	// Create the out and errOut channels
	out := make(chan string)
	errOut := make(chan error)

	// Call the readJSONObjects function in a separate goroutine
	go p.readJSONObjects(file.Content, out, errOut)

	// Loop until both channels are closed and set to nil
	for {
		select {
		case jsonStr, ok := <-out:
			if !ok {
				out = nil // Set the channel to nil to stop checking it
			} else {
				// Process the JSON string (e.g., pass it to Product.SetState())
				p.Logger.Info("Parser.Parse() Parsing JSON Object:", jsonStr)

				if file.IsProduct {
					p.Logger.Info("Parser.Parse() Parsing Product...")

					product := models.NewProduct(
						models.NewProductWithLogger(*modelsLogger),
					)
					var err error
					err = product.Parse(jsonStr)
					if err != nil {
						p.Logger.Warning("Parser.Parse() Failed to parse product while calling product.Parse", err)
					}
					err = product.Save()
					if err != nil {
						p.Logger.Error("Failed to save product", err)
					} else {
						productCount++ // Increment product count on successful save
					}
					file.IsProduct = true
					file.OwnerID = product.ID
					file.IsParsed = true
					file.LastModified = time.Now()

				} else {
					p.Logger.Info("Parsing Article...")
					article := models.NewArticle(
						models.NewArticleWithLogger(*modelsLogger),
					)
					err := article.Parse(jsonStr)
					if err != nil {
						p.Logger.Error("Parser.Parse() Failed to parse article", err)
					}
					err = article.Save()
					if err != nil {
						p.Logger.Error("Parser.Parse() Failed to save Article", err)
					} else {
						articleCount++ // Increment article count on successful save
					}
					file.IsProduct = false
					file.OwnerID = article.ID
					file.IsParsed = true
					file.LastModified = time.Now()
				}
				p.Logger.Info(fmt.Sprintf("Parser.Parse() Parsed %d Articles and %d Products", articleCount, productCount))
			}
		case err, ok := <-errOut:
			if !ok {
				errOut = nil // Set the channel to nil to stop checking it
			} else {
				// Handle the error (e.g., log it, exit the program, or take other appropriate action)
				p.Logger.Error("Parser.Parse() Error while parsing JSON Object", err)
				fmt.Println("Error:", err)
			}
		}

		// Break the loop when both channels are closed and set to nil
		if out == nil && errOut == nil {
			p.Logger.Info("Parser.Parse() Finished parsing file")
			break
		}
	}
	p.Logger.Info(fmt.Sprintf("Parsed a total of %d Articles and %d Products", articleCount, productCount))
}

func (p *Parser) readJSONObjects(content string, out chan<- string, errOut chan<- error) {
	p.Logger.Debug("readJSONObjects() go routine called")
	defer close(out)
	defer close(errOut)

	// Declare a variable to store each JSON object
	var jsonObj map[string]interface{}

	// Use strings.NewReader to read the content string
	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)

	// Set the buffer size to 50MB
	const bufferSize = 50 * 1024 * 1024 // 50MB buffer size
	buffer := make([]byte, bufferSize)
	scanner.Buffer(buffer, bufferSize)

	// Skip the first line
	if scanner.Scan() {
	} // Do nothing, just skip the first line

	// Iterate over each line in the file
	p.Logger.Info("readJSONObjects() Iterating over each line in the file")
	for scanner.Scan() {
		// Get the current line as a string
		line := scanner.Text()

		// Unmarshal the JSON line into the jsonObj variable
		err := json.Unmarshal([]byte(line), &jsonObj)
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

func (p *Parser) GetFiles() []string {
	retVal := []string{}
	dir := os.Getenv("DOWNLOAD_PATH")
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
