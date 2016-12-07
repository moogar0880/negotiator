package negotiator

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testContentNegotiatorType = "application/negotiated+json"
)

var invalidMediaType = errors.New("Invalid Media Type")

// testCN implements the ContentNegotiator interface for use in testing
type testCN struct {
	Foo string
	Bar int
}

func newTcn(foo string, bar int) *testCN {
	return &testCN{foo, bar}
}

func (tcn *testCN) ContentType(*Accept) (string, error) {
	return testContentNegotiatorType, nil
}

func (tcn *testCN) MarshalMedia(a *Accept) ([]byte, error) {
	switch a.MediaRange {
	case testContentNegotiatorType:
		data, err := json.Marshal(tcn)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, invalidMediaType
}

func (tcn *testCN) UnmarshalMedia(cType string, params ContentTypeParams, body []byte) error {
	switch cType {
	case testContentNegotiatorType:
		err := json.Unmarshal(body, &tcn)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestUnmarshalRequest(t *testing.T) {
	testIO := []struct {
		cType    string
		body     string
		err      error
		expected testCN
	}{
		// test simple JSON case
		{testContentNegotiatorType, `{"foo": "baz", "bar": 12}`, nil, *newTcn("baz", 12)},
		// test with no content type header set
		{"", `{"foo": "baz", "bar": 12}`, NoContentTypeErr, testCN{}},
		// test with invalid media type
		{"white space", `{"foo": "baz", "bar": 12}`,
			errors.New("mime: expected slash after first token"), testCN{}},
	}

	for _, test := range testIO {
		t.Run(test.cType, func(t *testing.T) {
			// create new http request
			req, _ := http.NewRequest("PUT", "http://example.com",
				bytes.NewReader([]byte(test.body)))
			req.Header[ContentTypeHeader] = []string{test.cType}

			// throw a testCN through with our request and ensure we see expected
			// results
			var tcn testCN
			res := UnmarshalMedia(req, &tcn)
			if res != nil && test.err != nil && res.Error() != test.err.Error() {
				t.Errorf("Expected Error %#v, got %#v instead", test.err, res)
			} else if res == nil && test.err != nil {
				t.Errorf("Expected error %v to be returned", test.err.Error())
			}

			// ensure everything unmarshalled correctly
			if tcn != test.expected {
				t.Errorf("Expected Result %#v, got %#v instead", test.expected, tcn)
			}
		})
	}
}

func TestMarshalMedia(t *testing.T) {
	testIO := []struct {
		inp        *testCN
		mediaRange mediaRange
		err        error
	}{
		// zero value content negotiatior
		{&testCN{}, testContentNegotiatorType, nil},
		// invalid media type
		{&testCN{}, "application/json", invalidMediaType},
	}

	for _, test := range testIO {
		t.Run(test.inp.Foo, func(t *testing.T) {
			writer := httptest.NewRecorder()
			err := MarshalMedia(writer,
				test.inp,
				&Accept{MediaRange: test.mediaRange})

			assert.Equal(t, test.err, err)
		})
	}
}
