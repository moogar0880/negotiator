package negotiator

import (
	"errors"
	"mime"
	"reflect"
)

var (
	// ErrNoContentType is the error returned if an accept header cannot be matched
	// in the current registry
	ErrNoContentType = errors.New("No Acceptable Content Type")
)

// Registry is a content type registry used for managing a mapping of media
// ranges to the interfaces that represent those resources
type Registry map[string]interface{}

// NewRegistry returns an empty Registry
func NewRegistry() *Registry {
	return &Registry{}
}

// Register registers the default struct value for a content type in the
// registry
func (r Registry) Register(contentType string, defaultValue interface{}) {
	if reflect.TypeOf(defaultValue).Kind() == reflect.Ptr {
		r[contentType] = reflect.ValueOf(defaultValue).Elem().Interface()
		return
	}
	r[contentType] = defaultValue
}

// Negotiate attempts to negotiate the proper interface for the provided accept
// header. Negotiate returns a copy of the default interface that best matches
// the provided accept header, if a match is found
func (r Registry) Negotiate(header string) (interface{}, error) {
	acceptHeader, err := ParseHeader(header)
	if err != nil {
		return nil, err
	}

	for _, hdr := range acceptHeader {
		if val, ok := r[string(hdr.MediaRange)]; ok {
			return reflect.ValueOf(val).Interface(), nil
		}
	}
	return nil, ErrNoContentType
}

// ContentType parses the provided Content-Type header and attempts to find an
// interface which implements the specified content type
func (r Registry) ContentType(header string) (interface{}, error) {
	mediaType, _, err := mime.ParseMediaType(header)
	if err != nil {
		return nil, err
	}

	if val, ok := r[mediaType]; ok {
		return reflect.ValueOf(val).Interface(), nil
	}
	return nil, ErrNoContentType
}