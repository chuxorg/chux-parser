package parsing

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chuxorg/chux-models/models"

	"github.com/chuxorg/chux-parser/s3"
)

// Parser struct for parsing
type Parser struct {
	products []models.Product
	articles []models.Article
}

// New returns a new Parser struct
func New(options ...func(*Parser)) *Parser {

	parser := &Parser{}
	for _, option := range options {
		option(parser)
	}
	return parser
}

func (p *Parser) Parse(file s3.File) {

	// Create the out and errOut channels
	out := make(chan string)
	errOut := make(chan error)

	// Call the readJSONObjects function in a separate goroutine
	go readJSONObjects(file.Content, out, errOut)

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
					product := models.NewProduct()
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
					article := models.NewArticle()
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

func readJSONObjects(content string, out chan<- string, errOut chan<- error) {

	defer close(out)
	defer close(errOut)

	// Declare a variable to store each JSON object
	var jsonObj map[string]interface{}

	// Use strings.NewReader to read the content string
	reader := strings.NewReader(content)

	// Create a new scanner to read the content line by line
	scanner := bufio.NewScanner(reader)

	// Skip the first line
	if scanner.Scan() {
	} // Do nothing, just skip the first line

	// Iterate over each line in the file
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
