package errors

import (
	"log"
)

// ChuxParserError is a custom error type
// that wraps an error and adds a message
// to the error.
// This is the error that is returned by
// all functions in chux-models that return
// an error.
type ChuxParserError struct {
	// Message is the message that is
	// given by chux-models when an error
	// occurs.
	// This message is used to provide
	// more context to the error.
	// The Err field contains the actual
	// error that occurred.
	Message  string
	InnerErr error
}

// NewChuxParserError returns a new ChuxModelsError
func NewChuxParserError(message string, err error) *ChuxParserError {
	return &ChuxParserError{
		Message:  message,
		InnerErr: err,
	}
}

func (e *ChuxParserError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error without
// the message added by chux-parser.
func (e *ChuxParserError) Unwrap() error {
	return e.InnerErr
}

// handleError is a helper function that handles
// errors occurring in chux-parser. This means
// that it prints the error message and the
// underlying error. It will also log the error
func handleError(err error) {
	log.Printf("Error: %v\n", err)
}
