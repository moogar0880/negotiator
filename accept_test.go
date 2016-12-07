package negotiator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMediaRange(t *testing.T) {
	testio := []struct {
		inp    mediaRange
		typ    string
		subtyp string
		suffix string
	}{
		{mediaRange("application/json"), "application", "json", ""},
		{mediaRange("application/*"), "application", "*", ""},
		{mediaRange("*/*"), "*", "*", ""},
		{mediaRange("application/json;indent=4"), "application", "json", ""},
		{mediaRange("application/resource+json;indent=4"), "application", "resource+json", "json"},
		{mediaRange("application resource"), "", "", ""},
	}

	for _, test := range testio {
		t.Run(test.inp.Type(), func(t *testing.T) {
			assert.Equal(t, test.typ, test.inp.Type(), "Types did not match")
			assert.Equal(t, test.subtyp, test.inp.SubType(), "Subtypes did not match")
			assert.Equal(t, test.suffix, test.inp.Suffix(), "Suffixes did not match")
		})
	}
}

func TestBadMediaRange(t *testing.T) {
	testio := []struct {
		inp string
		err error
	}{
		{"application/json", nil},
		{"application/*", nil},
		{"*/*", nil},
		{"application/json;indent=4", nil},
		{"application/resource+json;indent=4", nil},
		{"application resource", ErrInvalidMediaRange},
	}

	for _, test := range testio {
		t.Run(test.inp, func(t *testing.T) {
			if _, err := ParseAccept(test.inp); err != test.err {
				t.Errorf("Expected %s, got %s", test.err, err)
			}
		})
	}
}

func TestAcceptParams(t *testing.T) {
	testio := []struct {
		inp      string
		expected map[string]string
	}{
		{"application/json", map[string]string{}},
		{"application/json;indent=4",
			map[string]string{"indent": "4"}},
		{"application/json;indent=4; charset=utf8",
			map[string]string{"indent": "4", "charset": "utf8"}},
	}

	for _, test := range testio {
		t.Run(test.inp, func(t *testing.T) {
			acpt, err := ParseAccept(test.inp)

			assert.Nil(t, err, "Unable to parse valid header: %s", test.inp)

			for k, v := range test.expected {
				assert.Contains(t, acpt.AcceptParams, k,
					"Expected key %s not found in parsed params", k)
				assert.Equal(t, acpt.AcceptParams[k], v,
					"expected %s: %s to equal %s", k, acpt.AcceptParams[k], v)
			}
		})
	}
}

func TestBadParams(t *testing.T) {
	testio := []struct {
		inp  string
		fail bool
	}{
		{"application/json", false},
		{"application/json;q=0.3", false},
		{"application/json;foo=bar", false},
		{"application/json;foobar", true},
	}

	for _, test := range testio {
		t.Run(test.inp, func(t *testing.T) {
			_, err := ParseAccept(test.inp)

			failed := err == nil && test.fail == true
			assert.False(t, failed,
				"Expected header %s to contain a bad header param", test.inp)
		})
	}
}

func TestAcceptQuality(t *testing.T) {
	testio := []struct {
		inp      string
		expected float64
	}{
		{"application/json", 0.9},
		{"application/json;q=0.3", 0.3},
		{"application/json;indent=4", 1.0},
		{"application/json;indent=4;q=0.7", 0.7},
		{"application/json;indent=4; q=0.4", 0.4},
	}

	for _, test := range testio {
		t.Run(test.inp, func(t *testing.T) {
			acpt, err := ParseAccept(test.inp)
			assert.Nil(t, err, "Unable to parse valid header: %s", test.inp)
			assert.Equal(t, test.expected, acpt.Quality,
				"Expected quality of %f, got %f instead", test.expected, acpt.Quality)
		})
	}
}

func TestBadQuality(t *testing.T) {
	testio := []struct {
		inp  string
		fail bool
	}{
		{"application/json", false},
		{"application/json;q=0.3", false},
		{"application/json;q=1", false},
		{"application/json;q=foobar", true},
	}

	for _, test := range testio {
		t.Run(test.inp, func(t *testing.T) {
			_, err := ParseAccept(test.inp)

			failed := err == nil && test.fail == true
			assert.False(t, failed,
				"Expected header %s to contain a bad quality value", test.inp)
		})
	}
}

func TestAcceptExtensions(t *testing.T) {
	testio := []struct {
		inp      string
		expected map[string]string
	}{
		{"application/json", map[string]string{}},
		{"application/json;indent=4;q=1.0;version=1",
			map[string]string{"version": "1"}},
		{"application/json;indent=4; q=1.0; version=2",
			map[string]string{"version": "2"}},
	}

	for _, test := range testio {
		t.Run(test.inp, func(t *testing.T) {
			acpt, err := ParseAccept(test.inp)
			assert.Nil(t, err, "Unable to parse valid header: %s", test.inp)

			for k, v := range test.expected {
				assert.Contains(t, acpt.AcceptExt, k,
					"Expected key %s not found in parsed params", k)
				assert.Equal(t, acpt.AcceptExt[k], v,
					"expected %s: %s to equal %s", k, acpt.AcceptExt[k], v)
			}
		})
	}
}
