package errors

import "fmt"

// ChuxParserError is a custom error type
// that wraps an error and adds a message
// to the error.
// This is the error that is returned by
// all functions in chux-parse that return
// an error.
type ChuxParserError struct {
	// Message is the message that is
	// given by chux-parser when an error
	// occurs.
	// This message is used to provide
	// more context to the error.
	// The Err field contains the actual 
	// error that occurred.
	Message string
	Err     error
}

// NewChuxParserError returns a new ChuxParserError
func NewChuxParserError(message string, err error) *ChuxParserError {
	return &ChuxParserError{
		Message: message,
		Err:     err,
	}
}

func (e *ChuxParserError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error without
// the message added by chux-parser.
func (e *ChuxParserError) Unwrap() error {
	return e.Err
}


// handleError is a helper function that handles
// errors occuring in chux-parser. This means 
// that it prints the error message and the
// underlying error. It will also log the error	
func HandleError(err error) {
	fmt.Printf("Error: %v\n", err)
}