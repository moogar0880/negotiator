package negotiator

import (
	"testing"
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
		if test.typ != test.inp.Type() {
			t.Errorf("Expected Type %s, got %s instead", test.typ, test.inp.Type())
		}

		if test.subtyp != test.inp.SubType() {
			t.Errorf("Expected SubType %s, got %s instead", test.subtyp, test.inp.SubType())
		}

		if test.suffix != test.inp.Suffix() {
			t.Errorf("Expected Suffix %s, got %s instead", test.suffix, test.inp.Suffix())
		}
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
		if _, err := ParseAccept(test.inp); err != test.err {
			t.Errorf("Expected %s, got %s", test.err, err)
		}
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
		acpt, err := ParseAccept(test.inp)
		if err != nil {
			t.Errorf("Unable to parse valid header: %s", test.inp)
		}

		for k, v := range test.expected {
			if val, ok := acpt.AcceptParams[k]; ok {
				if v != val {
					t.Errorf("expected %s: %s to equal %s", k, val, v)
				}
			} else {
				t.Errorf("Expected key %s not found in parsed params", k)
			}
		}
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
		_, err := ParseAccept(test.inp)

		if err == nil && test.fail == true {
			t.Errorf("Expected header %s to contain a bad header param", test.inp)
		}
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
		acpt, err := ParseAccept(test.inp)
		if err != nil {
			t.Errorf("Unable to parse valid header: %s", test.inp)
		}

		if acpt.Quality != test.expected {
			t.Errorf("Expected quality of %f, got %f instead", test.expected, acpt.Quality)
		}
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
		_, err := ParseAccept(test.inp)

		if err == nil && test.fail == true {
			t.Errorf("Expected header %s to contain a bad quality value", test.inp)
		}
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
		acpt, err := ParseAccept(test.inp)
		if err != nil {
			t.Errorf("Unable to parse valid header: %s", test.inp)
		}

		for k, v := range test.expected {
			if val, ok := acpt.AcceptExt[k]; ok {
				if v != val {
					t.Errorf("expected %s: %s to equal %s", k, val, v)
				}
			} else {
				t.Errorf("Expected key %s not found in parsed params", k)
			}
		}
	}
}
