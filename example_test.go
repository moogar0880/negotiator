package negotiator_test

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/moogar0880/negotiator"
)

func ExampleRegistry() {
	type Message struct {
		Name     string
		Greeting string
	}

	type greeting struct {
		Phrase   string
		Language string
	}

	type MessageV2 struct {
		Name     string
		Greeting greeting
	}

	// define a registry that can handle the media types for our message structs
	registry := negotiator.NewRegistry()
	registry.Register("application/vnd.message.v1+json", Message{})
	registry.Register("application/vnd.message.v2+json", MessageV2{})

	// negotiate a predefined accept header, erroring if we don't support any of
	// the provided media types. Spoiler Alert - We do.
	acptHeader := "application/json, application/vnd.message.v1+json"
	model, accept, err := registry.Negotiate(acptHeader)
	if err != nil {
		fmt.Printf("Invalid Accept Header: %s", accept.MediaRange)
		return
	}

	// Dump our response (json encoded) into a buffer and print the result
	w := bytes.NewBuffer([]byte{})
	switch model.(type) {
	default:
		msg := MessageV2{Name: "John Doe",
			Greeting: greeting{
				Phrase:   "Hello",
				Language: "English"},
		}
		json.NewEncoder(w).Encode(&msg)
	case Message:
		msg := Message{Name: "John Doe", Greeting: "Hello"}
		json.NewEncoder(w).Encode(&msg)
	}
	fmt.Println(w.String())
	// Output:
	// {"Name":"John Doe","Greeting":"Hello"}
}
