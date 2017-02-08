package negotiator

import (
	"sort"
	"strings"
)

// AcceptHeader is a slice of individual Accept instances representing an
// entire accept header
type AcceptHeader []*Accept

// By is the type of a "less" function that defines the ordering of its Accept
// arguments.
type by func(a1, a2 *Accept) bool

// Sort is a method on the function type, By, that sorts the argument slice
// according to the function.
func (by by) Sort(accept AcceptHeader) {
	rs := &acceptSorter{
		accepts: accept,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(rs)
}

// recordSorter joins a By function and a slice of Records to be sorted
type acceptSorter struct {
	accepts AcceptHeader
	by      func(a1, a2 *Accept) bool // Closure used in the Less method
}

// Len is part of sort.Interface
func (a *acceptSorter) Len() int {
	return len(a.accepts)
}

// Swap is part of sort.Interface
func (a *acceptSorter) Swap(i, j int) {
	a.accepts[i], a.accepts[j] = a.accepts[j], a.accepts[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure
// in the sorter
func (a *acceptSorter) Less(i, j int) bool {
	return a.by(a.accepts[i], a.accepts[j])
}

// byWeight is a "by" closure which sorts based on an Accept's Quality field
func byWeight(a1, a2 *Accept) bool {
	return a1.Quality > a2.Quality
}

// ParseHeader parses an entire Accept header into an AcceptHeader instance and
// sorts it according to the relative quality of the accept headers provided
func ParseHeader(header string) (AcceptHeader, error) {
	var act *Accept
	var accepts AcceptHeader
	var err error

	values := strings.Split(header, ",")
	for _, value := range values {
		act = NewAccept()
		err = act.Parse(strings.TrimSpace(value))
		if err != nil {
			return nil, err
		}
		accepts = append(accepts, act)
	}

	by(byWeight).Sort(accepts)
	return accepts, nil
}
