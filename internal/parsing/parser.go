package parsing

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/chuxorg/chux-models/models/articles"
	"github.com/chuxorg/chux-models/models/products"
	"github.com/chuxorg/chux-parser/config"
	cfg "github.com/chuxorg/chux-parser/config"
)

// Parser struct for parsing
type Parser struct {
	products []products.Product
	articles []articles.Article
}

var _cfg *cfg.ParserConfig

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

func WithConfig(config  cfg.ParserConfig) func(*Parser) {
	return func(product *Parser) {
		_cfg = &config
	}
}

func (p *Parser) Parse(fileName string) {

	isProduct := true
    
	// Set the bizObjConfig variable to the value of the BizObjConfig field in the ParserConfig struct
	bizObjConfig := config.GetBizObjConfig(*_cfg)
	// Create the out and errOut channels
	out := make(chan string)
	errOut := make(chan error)

	// Call the readJSONObjects function in a separate goroutine
	go readJSONObjects(fileName, out, errOut)

	// Loop until both channels are closed and set to nil
	for {
		select {
		case jsonStr, ok := <-out:
			if !ok {
				out = nil // Set the channel to nil to stop checking it
			} else {
				// Process the JSON string (e.g., pass it to Product.SetState())
				fmt.Println("JSON Object:", jsonStr)
				if isProduct {			
					product := products.New(
						products.WithBizObjConfig(bizObjConfig),
					)
					product.SetState(jsonStr)
					product.Save()
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


// readJSONObjects reads a JSON file line by line, processes each JSON object,
// and sends the JSON object as a string to the output channel.
// It also sends any errors encountered to the error output channel.
func readJSONObjects(filePath string, out chan<- string, errOut chan<- error) {
	// Close both output channels after the function exits
	defer close(out)
	defer close(errOut)

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
	// Declare a variable to store each JSON object
	var jsonObj map[string]interface{}

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

