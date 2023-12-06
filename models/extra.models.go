package models

// CustomError represents a custom error with status and error message to handle error from database error with status.
type CustomError struct {
	Status int
	Err    error
}
