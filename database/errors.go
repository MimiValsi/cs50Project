package database

import (
	"errors"
)

// global variable to be used in database files
// Sends a message to dev to warning that a something went wrong
// with PSQL records
var (
	ErrNoRecord = errors.New("models: No matching record found")
)
