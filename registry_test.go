package negotiator

import (
	"reflect"
	"testing"
)

const (
	appJSON = "application/json"
	appXML  = "application/xml"
)

type testGeneric struct {
	X int
}

type testSpecific struct {
	X int
	Y int
}

func TestRegistryNegotiate(t *testing.T) {
	testReg := NewRegistry()
	testReg.Register("application/json", testGeneric{})
	testReg.Register("application/vnd.dyn.zone+json", &testSpecific{})
	testReg.Register("application/xhtml+xml", &testSpecific{})
	testReg.Register("*/*", &testGeneric{})

	testio := []struct {
		inp      string
		expected interface{}
		err      string
	}{
		{"application/json", testGeneric{}, ""},
		{"application/xml", nil, ErrNoContentType.Error()},
		{"application/json;foo", nil, ErrInvalidAcceptParam.Error()},
		{"*/*,application/json,application/vnd.dyn.zone+json;format=foo,application/*",
			testSpecific{}, ""},
		{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,/;q=0.8",
			testSpecific{}, ""},
		{"text/html, application/xhtml+xml, application/xml;q=0.9, image/webp, /;q=0.8",
			testSpecific{}, ""},
		{"image/jpeg, image/webp, */*",
			testGeneric{}, ""},
	}

	var i interface{}
	var err error
	for _, test := range testio {
		i, _, err = testReg.Negotiate(test.inp)

		if i != nil && test.expected == nil {
			t.Errorf("Expected %s to return nil. Got %+v instead", test.expected, i)
		}

		if i != test.expected {
			seen := reflect.TypeOf(i)
			expected := reflect.TypeOf(test.expected)
			t.Errorf("Expected type %s, got %s instead. Header (%s)", expected, seen, test.inp)
		}

		if len(test.err) > 0 && err != nil && test.err != err.Error() {
			t.Errorf("Expected error %q, got %q instead", test.err, err.Error())
		} else if len(test.err) > 0 && err == nil {
			t.Errorf("Expected an error, but got success from header %s", test.inp)
		}
	}
}

func TestRegistryContentType(t *testing.T) {
	testReg := NewRegistry()
	testReg.Register("application/json", testGeneric{})
	testReg.Register("application/vnd.dyn.zone+json", &testSpecific{})

	testio := []struct {
		inp      string
		expected interface{}
		err      string
	}{
		{"application/json", testGeneric{}, ""},
		{"application/xml", nil, ErrNoContentType.Error()},
		{"application/xml/json/ foobar", nil,
			"mime: unexpected content after media subtype"},
	}

	var i interface{}
	var err error
	for _, test := range testio {
		i, _, err = testReg.ContentType(test.inp)

		if i != nil && test.expected == nil {
			t.Errorf("Expected %s to return nil. Got %+v instead", test.expected, i)
		}

		if i != test.expected {
			seen := reflect.TypeOf(i)
			expected := reflect.TypeOf(test.expected)
			t.Errorf("Expected type %s, got %s instead. Header (%s)", expected, seen, test.inp)
		}

		if len(test.err) > 0 && err != nil && test.err != err.Error() {
			t.Errorf("Expected error %q, got %q instead", test.err, err.Error())
		} else if len(test.err) > 0 && err == nil {
			t.Errorf("Expected an error, but got success from header %s", test.inp)
		}
	}
}
