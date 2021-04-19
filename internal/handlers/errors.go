package handlers

// These errors are returned by the services and can be used to provide error codes to the
// API results.
const (
	ErrInvalidJSONInput HandlerError = "handlers: invalid_json, provided input cannot be parsed"
)

// PublicError is an error that returns a string code that can be presented to the API user.
type PublicError interface {
	error
	Code() string
	Detail() string
}

// HandlerError defines errors exported by this package. This type implement a Code() method that
// extracts a unique error code defined for each error value exported.
type HandlerError string

// Error returns the exact original message of the e value.
func (e HandlerError) Error() string {
	return string(e)
}

// Code extracts the error code string present on the value of e.
//
// An error code is defined as the string after the package prefix and colon, and before the comma that follows this string. Example:
//		"handlers: error_code, this is a validation error"
func (e HandlerError) Code() string {
	// remove the prefix
	s := string(e)[len("handlers: "):]

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
//		"handlers: error_code, this is the error detail string"
func (e HandlerError) Detail() string {
	// remove the prefix
	s := string(e)[len("handlers: "):]

	// extract the error code
	for i := 1; i < len(s); i++ {
		if s[i] == ',' {
			s = s[i+2:] // +2 removes the comma and the space
			break
		}
	}

	return s
}
