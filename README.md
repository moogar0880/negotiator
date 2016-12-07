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

## Example Usage
```go
import (
  "encoding/json"
  "net/http"

  "github.com/moogar0880/negotiator"
)

var Registry negotiator.Registry

type Message struct {
  Name string
  Greeting string
}

type greeting struct {
  Phrase string
  Language string
}

type MessageV2 struct {
  Name string
  Greeting greeting
}

func init() {
  Registry = negotiator.NewRegistry()
  Registry.Register("application/vnd.message.v1+json", Message{})
  Registry.Register("application/vnd.message.v2+json", MessageV2{})
}

func MessageHandler(w http.ResponseWriter, req *http.Request) {
  model, accept, err := Registry.Negotiate(req.Header.Get("Accept"))
  if err != nil {
    http.Error(w, "Invalid Accept Header", http.StatusNotAcceptable)
    return
  }

  switch model.(type){
  default:
    msg := MessageV2{Name: "John Doe",
      Greeting: greeting{
        Phrase: "Hello",
        Language: "English"},
    }
    json.NewEncoder(w).Encode(msg)
  case Message:
    msg := Message{Name: "John Doe", Greeting: "Hello"}
    json.NewEncoder(w).Encode(msg)
  }
}
```
