package schema

import "errors"

var (
	ErrNoDefs        = errors.New("no definition of variable in schema")
	ErrNoResolver    = errors.New("no resolver available for the tag")
	ErrInvalidScheme = errors.New("invalid scheme")
)

// IData simply gives the data corresponding to column name.
type IData interface {
	Get(string) interface{}
}

// IResolver resolves the variable of given name to the Token and
// resolves the conditions when the database entry match with the
// input entry.
type IResolver interface {
	AddResolver(tag string, resolver func(string) (interface{}, error))

	Resolve(name string, value string) (Token, error)

	Equals(name string, v1, v2 interface{}) (bool, error)

	AddEqualChecker(name string, checker func(v1, v2 interface{}) (bool, error))

	Columns() []string
}
