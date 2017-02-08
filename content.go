package negotiator

import (
	"io"
	"io/ioutil"
	"mime"
	"net/http"
)

// ContentTypeHeader is the constant value for the key indicating the
// Content-Type header
const ContentTypeHeader = "Content-Type"

// The ContentNegotiator interface defines the mechanism through which arbitrary
// interfaces can be provided information about the provided Accept and
// Content-Type headers, to control marshalling and unmarshalling
// request/response data as correctly as possible. Optionally, requests may be
// rejected if provided arguments are invalid, unacceptable, or otherwise
// erronenous for a given resource.
type ContentNegotiator interface {
	// ContentType accepts the provided Accept header struct and returns the
	// matched content type, or an error
	ContentType(*Accept) (string, error)

	// MarshalMedia returns a raw byte slice containing an appropriately rendered
	// representation of the provided resource, or an error.
	MarshalMedia(*Accept) ([]byte, error)

	// UnmarshalMedia accepts the content type and content type parameters
	// provided in a request, as well as the raw request body, and unmarshals it
	// into the ContentNegotiator implementation struct
	UnmarshalMedia(string, ContentTypeParams, []byte) error
}

// MarshalMedia marshals the ContentNegotiator to the provided io.Writer, based
// on an Accept. An error is returned if the ContentNegotiator's MarshalMedia
// call fails, or if the data can't be written to the io.Writer
func MarshalMedia(w io.Writer, cn ContentNegotiator, acpt *Accept) error {
	data, err := cn.MarshalMedia(acpt)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// UnmarshalMedia handles unmarshalling an http.Request body, using a
// ContentNegotiator instance. An error is returned if no Content-Type header
// was provided, if the provided Content-Type header was poorly formatted, or
// if the body of the http.Request could not be read.
func UnmarshalMedia(req *http.Request, cn ContentNegotiator) error {
	var header string
	if header = req.Header.Get(ContentTypeHeader); len(header) == 0 {
		return ErrNoContentType
	}

	mediaType, params, err := mime.ParseMediaType(header)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	cn.UnmarshalMedia(mediaType, params, body)
	return nil
}
