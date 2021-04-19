package models

// These errors are returned by the services and can be used to provide error codes to the
// API results.
const (
	ErrNotFound         ModelError = "models: not_found, resource not found"
	ErrInvalidJSONInput ModelError = "models: invalid_json, provided input cannot be parsed"
)

// CodeError is an error that returns a string code that can be presented to the API user.
type PublicError interface {
	error
	Code() string
	Detail() string
}

// ModelError defines errors exported by this package. This type implement a Code() method that
// extracts a unique error code defined for each error value exported.
type ModelError string

// Error returns the exact original message of the e value.
func (e ModelError) Error() string {
	return string(e)
}

// Code extracts the error code string present on the value of e.
//
// An error code is defined as the string after the package prefix and colon, and before the comma that follows this string. Example:
//		"models: error_code, this is a validation error"
func (e ModelError) Code() string {
	// remove the prefix
	s := string(e)[len("models: "):]

	// extract the error code
	for i := 1; i < len(s); i++ {
		if s[i] == ',' {
			s = s[:i]
			break
		}
	}

	return s
}

// Detail extracts the error detail string present on the value of e.
//
// An error detail is defined as the string after the package prefix and colon, and after the comma that follows this string. Example:
//		"models: error_code, this is the error detail string"
func (e ModelError) Detail() string {
	// remove the prefix
	s := string(e)[len("models: "):]

	// extract the error code
	for i := 1; i < len(s); i++ {
		if s[i] == ',' {
			s = s[i+2:] // +2 removes the comma and the space
			break
		}
	}

	return s
}
