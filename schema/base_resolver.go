package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// BaseResolver is a basic template to create resolvers from
type BaseResolver struct {
	sch       map[string]string
	resolvers map[string]func(string) (interface{}, error)
	checkers  map[string]func(v1, v2 interface{}) (bool, error)
	cols      []string
}

// Equals looks for the corresponding equality checker for the
// given name and uses it on the given values.
// *** If no checker for the given name is available and checker ***
// ***   for empty name ("") is available, Equals will use it.   ***
func (s BaseResolver) Equals(name string, v1, v2 interface{}) (bool, error) {
	// If no checker are available, use the universal checker at ""
	if f, ok := s.checkers[name]; ok {
		return f(v1, v2)
	} else if f, ok := s.checkers[""]; ok {
		return f(v1, v2)
	}
	return false, nil
}

func (s BaseResolver) Columns() []string {
	return s.cols
}

func (s *BaseResolver) AddResolver(tag string, resolver func(string) (interface{}, error)) {
	s.resolvers[tag] = resolver
}

func (s *BaseResolver) AddEqualChecker(name string, checker func(v1, v2 interface{}) (bool, error)) {
	s.checkers[name] = checker
}

func (s BaseResolver) Resolve(name string, value string) (Token, error) {
	tag, ok := s.sch[name]
	if !ok {
		return Token{}, ErrNoDefs
	}

	resolver, ok := s.resolvers[tag]
	if !ok {
		return Token{}, ErrNoResolver
	}

	val, err := resolver(value)
	if err != nil {
		return Token{}, err
	}

	return Token{
		Id:    name,
		Type:  tag,
		Value: val,
	}, nil
}

// For printing purposes
func (s BaseResolver) String() string {
	bts, err := json.MarshalIndent(s.sch, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(bts)
}

// Expects a byte stream of json
func ParseJSONToBaseResolver(r io.Reader) (*BaseResolver, error) {
	s := &BaseResolver{
		sch:       make(map[string]string),
		resolvers: make(map[string]func(string) (interface{}, error)),
		checkers:  make(map[string]func(v1, v2 interface{}) (bool, error)),
	}

	if err := json.NewDecoder(r).Decode(&s.sch); err != nil {
		return nil, errors.New(fmt.Sprintf("%s:%s", ErrInvalidScheme, err))
	}
	// Memoize the columns of the schema
	var keys []string
	for k := range s.sch {
		keys = append(keys, k)
	}
	s.cols = keys

	return s, nil
}
