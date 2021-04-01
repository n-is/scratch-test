package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"scratch-test/schema"
	"strings"
	"testing"
)

type dummySchema struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (d dummySchema) AddEqualChecker(string, func(interface{}, interface{}) (bool, error)) {
	panic("implement me")
}

func (d dummySchema) Get(s string) interface{} {
	switch s {
	case "name":
		return d.Name
	case "age":
		return d.Age
	}
	return nil
}

func (d dummySchema) AddResolver(_ string, _ func(string) (interface{}, error)) {
	panic("implement me")
}

func (d dummySchema) Resolve(_ string, _ string) (schema.Token, error) {
	panic("implement me")
}

func (d dummySchema) Equals(_ string, v1, v2 interface{}) (bool, error) {
	return strings.EqualFold(fmt.Sprintf("%v", v1), fmt.Sprintf("%v", v2)), nil
}

func (d dummySchema) Columns() []string {
	return []string{"name", "age"}
}

func TestMemoryDB(t *testing.T) {

	sch := dummySchema{}
	data := []dummySchema{
		{
			Name: "Good Health Home",
			Age:  55,
		},
		{
			Name: "Mayo Clinic",
			Age:  55,
		},
		{
			Name: "Good Health Home",
			Age:  50,
		},
		{
			Name: "Hopkins Hospital Baltimore",
			Age:  55,
		},
		{
			Name: "Mount Sinai Hospital",
			Age:  5342535,
		},
		{
			Name: "Tufts Medical Center",
			Age:  550,
		},
		{
			Name: "UAB Hospital",
			Age:  5,
		},
		{
			Name: "Swedish Medical Center",
			Age:  55,
		},
		{
			Name: "Scratchpay Test Pet Medical Center",
			Age:  55,
		},
	}

	database := CreateMemoryDB()
	err := database.Init(sch)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	// Enter the data twice
	for i := 0; i < 2; i++ {
		for _, d := range data {
			err := database.Insert(d)
			if err != nil {
				t.Errorf("Error: %s", err.Error())
			}
		}
	}

	entries, err := database.Read(map[string][]interface{}{
		"name": {"Good Health home"},
	}, -1)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	if len(entries) != 4 {
		t.Errorf("Error: Length Mismatch. Want: %d, Got: %d", 2, len(entries))
	}

	for _, d := range data {
		values, err := database.Read(map[string][]interface{}{
			"name": {d.Name},
			"age":  {d.Age},
		}, -1)
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
		if len(values) != 2 {
			t.Errorf("Error: Length Mismatch. Want: %d, Got: %d", 2, len(entries))
		}

		bts, err := json.Marshal(values[0])
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		var dummy dummySchema
		if err := json.NewDecoder(bytes.NewReader(bts)).Decode(&dummy); err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		if !reflect.DeepEqual(dummy, d) {
			t.Errorf("Error:  Want: %v, Got: %v", d, dummy)
		}
	}
}
