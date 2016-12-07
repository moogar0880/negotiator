package negotiator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func acceptWithMedia(media mediaRange, qual float64) *Accept {
	a := NewAccept()
	a.MediaRange = media
	a.Quality = qual
	return a
}

func TestParseHeader(t *testing.T) {
	testIO := []struct {
		inp    string
		expect AcceptHeader
		err    error
	}{
		// simple accept header, application/json
		{"application/json",
			AcceptHeader{
				acceptWithMedia("application/json", 0.9),
			}, nil},
		// default accept header for Google Chrome
		{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,/;q=0.8",
			AcceptHeader{
				acceptWithMedia("text/html", 0.9),
				acceptWithMedia("application/xhtml+xml", 0.9),
				acceptWithMedia("application/xml", 0.9),
				acceptWithMedia("image/webp", 0.9),
				acceptWithMedia("/", 0.8),
			}, nil},
		// default accept header for Google Chrome, with whitespace
		{"text/html, application/xhtml+xml, application/xml;q=0.9, image/webp, /;q=0.8",
			AcceptHeader{
				acceptWithMedia("text/html", 0.9),
				acceptWithMedia("application/xhtml+xml", 0.9),
				acceptWithMedia("application/xml", 0.9),
				acceptWithMedia("image/webp", 0.9),
				acceptWithMedia("/", 0.8),
			}, nil},
	}

	for _, test := range testIO {
		t.Run(test.inp, func(t *testing.T) {
			header, err := ParseHeader(test.inp)
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.expect, header)
		})
	}
}
