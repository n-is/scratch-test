package db

import (
	"encoding/json"
	"scratch-test/schema"
	"sync"
)

// Currently on equals filter is supported

type MemoryDB struct {
	resolver schema.IResolver
	mu       sync.RWMutex
	data     map[string][]interface{}
}

func (m *MemoryDB) Close() error {
	return nil
}

func (m *MemoryDB) Init(resolver schema.IResolver) error {
	m.resolver = resolver
	return nil
}

func (m *MemoryDB) Insert(d schema.IData) error {
	return m.Update(func(tx ITx) error {
		tx.Insert(d)
		return nil
	})

}

func (m *MemoryDB) Read(condition map[string][]interface{}, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := m.View(func(tx ITx) error {
		var err error
		results, err = tx.Read(condition, limit)
		return err
	})

	return results, err
}

func (m *MemoryDB) begin(writable bool) (ITx, error) {
	tx := &MemoryTx{
		db:       m,
		writable: writable,
	}
	tx.lock()

	return tx, nil
}

func (m *MemoryDB) managed(writable bool, fn func(tx ITx) error) error {
	tx, err := m.begin(writable)
	if err != nil {
		return err
	}
	defer tx.unlock()

	return fn(tx)
}

func (m *MemoryDB) View(fn func(tx ITx) error) error {
	return m.managed(false, fn)
}

func (m *MemoryDB) Update(fn func(tx ITx) error) error {
	return m.managed(true, fn)
}

func CreateMemoryDB() *MemoryDB {
	return &MemoryDB{data: make(map[string][]interface{})}
}

func (m *MemoryDB) String() string {
	// Print at max top 10 entries

	var results []map[string]interface{}

	err := m.View(func(tx ITx) error {
		for i := 0; i < 10; i++ {
			result, err := tx.ReadRow(i)
			if err == ErrNoEntry {
				break
			}
			results = append(results, result)
		}
		return nil
	})

	bts, err := json.MarshalIndent(results, "", "\t")
	if err != nil {
		return err.Error()
	}

	return string(bts)
}
