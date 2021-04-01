package db

import (
	"errors"
	"log"
	"scratch-test/schema"
)

var (
	ErrNoEntry = errors.New("no entry")
)

// MemoryTx handles the transaction for in-memory database
type MemoryTx struct {
	db       *MemoryDB
	writable bool
}

// Insert is just amortized constant time (for fixed schema)
func (tx *MemoryTx) Insert(d schema.IData) {
	cols := tx.db.resolver.Columns()
	for _, c := range cols {
		val := d.Get(c)
		tx.db.data[c] = append(tx.db.data[c], val)
	}
}

// TODO:
// Implement efficient Read. This has O(n) complexity, which is very very bad.
// Indexing is probably the best way.
func (tx *MemoryTx) Read(condition map[string][]interface{}, limit int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}

	cols := tx.db.resolver.Columns()
	if len(cols) == 0 {
		return nil, ErrNoEntry
	}

	if condition == nil {
		// Send all the data
		for i := range tx.db.data[cols[0]] {

			if limit > 0 && len(data) >= limit {
				break
			}

			d := make(map[string]interface{})
			for _, c := range cols {
				d[c] = tx.db.data[c][i]
			}

			data = append(data, d)
		}
	}

	// No need to check if no conditions are available
	if len(condition) == 0 {
		return data, nil
	}

	for i := range tx.db.data[cols[0]] {

		if limit > 0 && len(data) >= limit {
			break
		}

		d := make(map[string]interface{})
		for _, c := range cols {
			d[c] = tx.db.data[c][i]
		}

		pick := false
		first := true
		for key, val := range d {
			if values, ok := condition[key]; ok {
				if first {
					pick = true
					first = false
				}
				or := false
				for _, vl := range values {
					ok, err := tx.db.resolver.Equals(key, vl, val)
					if err != nil {
						log.Println(err)
						continue
					}
					if ok {
						or = true
						break
					}
				}
				pick = pick && or
			}
			if !pick && !first {
				break
			}
		}
		if pick {
			data = append(data, d)
		}
	}

	return data, nil
}

func (tx *MemoryTx) ReadRow(i int) (map[string]interface{}, error) {

	d := make(map[string]interface{})
	for _, c := range tx.db.resolver.Columns() {
		if len(tx.db.data[c]) > i {
			d[c] = tx.db.data[c][i]
		} else {
			return nil, ErrNoEntry
		}
	}

	return d, nil
}

func (tx *MemoryTx) lock() {
	if tx.writable {
		tx.db.mu.Lock()
	} else {
		tx.db.mu.RLock()
	}
}

func (tx *MemoryTx) unlock() {
	if tx.writable {
		tx.db.mu.Unlock()
	} else {
		tx.db.mu.RUnlock()
	}
}
