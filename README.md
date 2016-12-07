# negotiator
Negotiator is a content negotiation library aimed to support strong content
typing in RESTful HTTP services. This library implements both `Accept` and
`Content-Type` header parsers and struct variants that are fully compliant with
both RFC-6839 and RFC-7231.

[![Build Status](https://travis-ci.org/moogar0880/negotiator.svg?branch=master)](https://travis-ci.org/moogar0880/negotiator)
[![Go Report Card](https://goreportcard.com/badge/github.com/moogar0880/negotiator)](https://goreportcard.com/report/github.com/moogar0880/negotiator)
[![GoDoc](https://godoc.org/github.com/moogar0880/negotiator?status.svg)](https://godoc.org/github.com/moogar0880/negotiator)

## Installation
```bash
$ go get github.com/moogar0880/negotiator
```

## Simple Example
This simple example shows how to represent a basic Message resource using the
content type registry. It has a single, defined, media type, and a default
representation.

```go
import (
  "encoding/json"
  "encoding/xml"
  "net/http"

  "github.com/moogar0880/negotiator"
)

const v1JSONMediaType = "application/vnd.message.v1+json"

var (
  Registry *negotiator.Registry
  defaultMessage = Message{Greeting: "Hello", Name: "World"}
)

type Message struct {
  Name string
  Greeting string
}

// Return a content type for the Message resource type
func (m *Message) ContentType(a *Accept) (string, error) {
  return v1JSONMediaType, nil
}

// handle encoding a message resource into a byte slice
func (m *Message) MarshalMedia(a *Accept) ([]byte, error) {
  data, _ := json.Marshal(m)
  return data, nil
}

// handle unmarshalling a message from an http request body
func (m *Message) UnmarshalMedia(cType string, params ContentTypeParams, body []byte) error {
  json.Unmarshal(body, &tcn)
  return nil
}

func init() {
  Registry = negotiator.NewRegistry()
  Registry.Register("application/vnd.message.v1+json", defaultMessage)
}

func messageHandler(w http.ResponseWriter, req *http.Request) {
  model, accept, err := Registry.Negotiate(req.Header.Get("Accept"))
  if err != nil {
    http.Error(w, "Invalid Accept Header", http.StatusNotAcceptable)
    return
  }

  negotiator.MarshalMedia(accept, model)
}
```


## Example Usage With Versioned Resource
This example shows how to handle versioning resources at the media type level,
similar to how the [Github API](https://developer.github.com/v3/#current-version)
does.

Expanding from our previous example, let's introduce a `vnd.message.v2+json`
resource that contains information about what language the message Greeting
is in.

```go
import (
  "encoding/json"
  "net/http"

  "github.com/moogar0880/negotiator"
)

const (
  v1JSONMediaType = "application/vnd.message.v1+json"
  v2JSONMediaType = "application/vnd.message.v2+json"
)

var Registry negotiator.Registry

// the original v1 message resource
type MessageV1 struct {
  Name string
  Greeting string
}

type greeting struct {
  Phrase string
  Language string
}

// The new message resource with a Greeting object instead of a string
type Message struct {
  Name string
  Greeting greeting
}

// Return a content type matching what was requested in the accept header
func (m *Message) ContentType(a *Accept) (string, error) {
  switch a.MediaRange {
  case v1JSONMediaType:
    return v1JSONMediaType, nil
  case v2JSONMediaType:
    return v2JSONMediaType, nil
  }
  return "", errors.New("Unsupported Media Type")
}

// handle marshlling a message resource, including converting to the old
// message resource format, if that's what was requested
func (m *Message) MarshalMedia(a *Accept) ([]byte, error) {
  switch a.MediaRange {
  case v1JSONMediaType:
    data, _ := json.Marshal(MessageV1{Name: m.Name, Greeting: m.Greeting.Phrase})
    return data, nil
  case v2JSONMediaType:
    data, _ := json.Marshal(m)
    return data, nil
  }
  return nil, errors.New("Unsupported Media Type")
}

// handle unmarshalling messages in either format from an HTTP request body as
// would be seen with a POST or PUT
func (m *Message) UnmarshalMedia(cType string, params ContentTypeParams, body []byte) error {
	switch cType {
  case v1JSONMediaType:
    var m1 MessageV1
    json.Unmarshal(body, &m1)
    m.Name = m1.Name
    m.Greeting.Phrase = m1.Greeting
    return nil
	case v2JSONMediaType:
		json.Unmarshal(body, &tcn)
    return nil
	}
	return errors.New("Unsupported Media Type")
}

func init() {
  Registry = negotiator.NewRegistry()
  Registry.Register("application/vnd.message.v1+json", Message{})
  Registry.Register("application/vnd.message.v2+json", Message{})
}

func MessageHandler(w http.ResponseWriter, req *http.Request) {
  model, accept, err := Registry.Negotiate(req.Header.Get("Accept"))
  if err != nil {
    http.Error(w, "Invalid Accept Header", http.StatusNotAcceptable)
    return
  }

  MarshalMedia(w, model, accept)
}
```
