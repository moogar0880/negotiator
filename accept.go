package negotiator

import (
	"errors"
	"strconv"
	"strings"
)

const (
	// AcceptParamsQuality is the default quality of a media range with specified
	// accept-params. (e.g text/html;level=1)
	AcceptParamsQuality float64 = 1.0

	// MediaRangeSubTypeQuality is the default quality of a media range with
	// type and subtype defined. (e.g text/html)
	MediaRangeSubTypeQuality float64 = 0.9

	// MediaRangeWildcardSubtypeQuality is the default weight of a media range
	// with a wildcarded subtype. (e.g text/*)
	MediaRangeWildcardSubtypeQuality float64 = 0.8

	// MediaRangeWildcardQuality is the default quality for a wildcarded media
	// range. (e.g */*)
	MediaRangeWildcardQuality float64 = 0.7

	// WildCard is the constant character "*", representing an accept header
	// wildcard character
	WildCard string = "*"
)

var (
	// ErrInvalidMediaRange is the error returned when an invalid accept
	// media range is parsed
	ErrInvalidMediaRange = errors.New("Invalid Accept Media Range")

	// ErrInvalidAcceptParam is the error returned when an invalid accept
	// parameter is parsed
	ErrInvalidAcceptParam = errors.New("Invalid Accept Parameter")
)

type (
	// mediaRange is a string alias for an accept header's media range
	mediaRange string

	// mediaParams are a mapping of parameter to value argument strings. parsed
	// from arguments of the form "foo=bar"
	mediaParams map[string]string

	// acceptExt is a type alias to mediaParams for parameters that follow the
	// quality, "q", parameter in an accept header
	acceptExt mediaParams
)

// Type returns the type of the media range instance. eg, a media range
// "application/json" has a Type of "application"
func (m mediaRange) Type() string {
	var s = string(m)
	var idx = strings.Index(s, "/")
	if idx == -1 {
		return ""
	}
	return s[:idx]
}

// SubType returns the sub-type of the media range instance. eg, a media range
// "application/json" has a SubType of "json"
func (m mediaRange) SubType() string {
	var s = string(m)
	var start = strings.Index(s, "/")
	var end = strings.Index(s, ";")

	if start == -1 {
		return ""
	}

	if end < start {
		return s[start+1:]
	}
	return s[start+1 : end]
}

// Return a structured media type suffix, as defined by RFC-6839, if one is
// present. Eg, a media range of application/resource+json returns "json".
// Because a these suffixes are optional, if no suffix is present an empty
// string is returned
func (m mediaRange) Suffix() string {
	var s = string(m)
	var end = strings.Index(s, ";")
	var idx = strings.LastIndex(s, "+")
	if idx == -1 {
		return ""
	}
	return s[idx+1 : end]
}

// Accept is the struct representation of a single accept header value
type Accept struct {
	MediaRange   mediaRange
	AcceptParams mediaParams
	Quality      float64
	AcceptExt    acceptExt
}

// NewAccept returns a zero-valued Accept instance
func NewAccept() *Accept {
	return &Accept{AcceptParams: make(mediaParams),
		Quality:   -1.0,
		AcceptExt: make(acceptExt)}
}

// calculateQuality calculates the default quality of this accept value, if one
// was not explicitly provided
func (a *Accept) calculateQuality() {
	if a.Quality != -1.0 {
		return
	} else if len(a.AcceptParams) > 0 {
		a.Quality = AcceptParamsQuality
	} else if a.MediaRange.SubType() != WildCard {
		a.Quality = MediaRangeSubTypeQuality
	} else if a.MediaRange.Type() != WildCard && a.MediaRange.SubType() == WildCard {
		a.Quality = MediaRangeWildcardSubtypeQuality
	} else if a.MediaRange.Type() == WildCard && a.MediaRange.SubType() == WildCard {
		a.Quality = MediaRangeWildcardQuality
	}
}

// Parse parses the provided string argument into an Accept instance, returning
// an error if the provided value is not properly formatted
func (a *Accept) Parse(accept string) error {
	var idx = strings.Index(accept, "/")
	if idx == -1 {
		return ErrInvalidMediaRange
	}

	idx = strings.Index(accept, ";")
	a.parseMediaRange(accept, idx)

	// if there was more than a simple media range provided, parse it and return
	// any parsing errors that occur
	var err error
	if idx != -1 {
		err = a.parseAcceptParams(accept, idx)
	}

	a.calculateQuality()
	return err
}

// parseMediaRange parses the media range out of an accept header up to an
// optional ';' character
func (a *Accept) parseMediaRange(accept string, semiIndex int) {
	if semiIndex == -1 {
		a.MediaRange = mediaRange(accept)
	} else {
		a.MediaRange = mediaRange(accept[:semiIndex])
	}
}

// parseMediaRange parses the optional accept parameters, up to an optional "q"
// "quality" parameter, any the following accept extension parameters
func (a *Accept) parseAcceptParams(accept string, semiIndex int) error {
	var keyVal []string
	var qParsed bool
	var err error
	for _, param := range strings.Split(accept[semiIndex+1:], ";") {
		keyVal = strings.Split(strings.TrimSpace(param), "=")
		if len(keyVal) != 2 {
			return ErrInvalidAcceptParam
		}

		if keyVal[0] == "q" {
			err = a.parseQuality(keyVal[1])
			if err != nil {
				return err
			}
			qParsed = true
		} else if qParsed {
			a.AcceptExt[keyVal[0]] = keyVal[1]
		} else {
			a.AcceptParams[keyVal[0]] = keyVal[1]
		}
	}
	return nil
}

// parseQuality parses the value of a quality ("q") parameter value as a float
func (a *Accept) parseQuality(val string) error {
	flt, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}

	a.Quality = flt
	return nil
}

// ParseAccept parses the provided accept header and returns a newly created
// Accept struct, and a conditional error
func ParseAccept(header string) (*Accept, error) {
	acpt := NewAccept()
	err := acpt.Parse(header)
	if err != nil {
		return nil, err
	}
	return acpt, nil
}
