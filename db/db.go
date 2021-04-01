package db

import (
	"scratch-test/schema"
)

// IDatabase defines the interface for storing and reading
// the data. This allows us to switch from memory-backed
// database to disk-backed database & vice-versa easily.
type IDatabase interface {
	Init(schema.IResolver) error

	// For simple insertion & deletion. No race condition allowed.
	Insert(schema.IData) error

	Read(map[string][]interface{}, int) ([]map[string]interface{}, error)

	// For Transaction Driven Implementation.
	begin(writable bool) (ITx, error)

	View(fn func(tx ITx) error) error

	Update(fn func(tx ITx) error) error

	Close() error
}

// ITx defines the interface for the transactions. Transactions
// helps to support atomic operations on internal database.
type ITx interface {
	Insert(d schema.IData)

	Read(map[string][]interface{}, int) ([]map[string]interface{}, error)

	ReadRow(i int) (map[string]interface{}, error)

	lock()

	unlock()
}
